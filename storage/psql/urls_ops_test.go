package psql

import (
	"errors"

	"github.com/google/uuid"

	"github.com/vstdy0/go-shortener/model"
	"github.com/vstdy0/go-shortener/pkg"
)

func (s *TestSuite) TestURLs_HasURL() {
	s.Run("Has non-existing url", func() {
		res, err := s.storage.HasURL(s.ctx, 10)
		s.Require().NoError(err)
		s.Assert().False(res)
	})

	s.Run("Has url_1", func() {
		expectedURL := s.fixtures.URLS[0].ToCanonical()

		res, err := s.storage.HasURL(s.ctx, expectedURL.ID)
		s.Require().NoError(err)
		s.Assert().True(res)
	})
}

func (s *TestSuite) TestURLs_AddURLs() {
	urlsToAdd := []model.URL{
		{
			UserID: uuid.New(),
			URL:    "https://extremely-lengthy-url-3.com/",
		},
		{
			UserID: uuid.New(),
			URL:    "https://extremely-lengthy-url-4.com/",
		},
	}

	s.Run("Add non-existing urls", func() {
		res, err := s.storage.AddURLs(s.ctx, urlsToAdd)
		s.Require().NoError(err)

		for idx := range urlsToAdd {
			s.Assert().EqualValues(urlsToAdd[idx].UserID, res[idx].UserID)
			s.Assert().EqualValues(urlsToAdd[idx].URL, res[idx].URL)
			s.Assert().NotEqual(0, res[idx].ID)
		}
	})

	s.Run("Add existing url", func() {
		existingURL, err := s.fixtures.URLS.ToCanonical()
		s.Require().NoError(err)

		res, err := s.storage.AddURLs(s.ctx, existingURL)
		s.Require().Error(err)
		s.Require().True(errors.Is(err, pkg.ErrAlreadyExists))
		s.Assert().EqualValues(existingURL[0].ID, res[0].ID)
	})
}

func (s *TestSuite) TestURLs_GetURL() {
	s.Run("Get non-existing url", func() {
		res, err := s.storage.GetURL(s.ctx, 10)
		s.Require().NoError(err)
		s.Require().EqualValues(model.URL{}, res)
	})

	s.Run("Get existing url", func() {
		expectedURL := s.fixtures.URLS[0].ToCanonical()

		res, err := s.storage.GetURL(s.ctx, expectedURL.ID)
		s.Require().NoError(err)
		s.Assert().EqualValues(expectedURL, res)
	})
}

func (s *TestSuite) TestURLs_GetUserURLs() {
	s.Run("Get non-existing user urls", func() {
		res, err := s.storage.GetUserURLs(s.ctx, uuid.New())
		s.Require().NoError(err)
		s.Require().Nil(res)
	})

	s.Run("Get existing user urls", func() {
		existingURLs, err := s.fixtures.URLS.ToCanonical()
		s.Require().NoError(err)

		res, err := s.storage.GetUserURLs(s.ctx, existingURLs[0].UserID)
		s.Require().NoError(err)
		s.Assert().EqualValues(existingURLs, res)
	})
}

func (s *TestSuite) TestURLs_RemoveUserURLs() {
	userURLs := []model.URL{
		{
			ID:     1,
			UserID: s.fixtures.URLS[0].UserID,
		},
	}
	foreignURLs := []model.URL{
		{
			ID:     2,
			UserID: uuid.New(),
		},
	}

	s.Run("Remove urls by creator", func() {
		err := s.storage.RemoveUserURLs(s.ctx, userURLs)
		s.Require().NoError(err)

		res, err := s.storage.GetURL(s.ctx, 1)
		s.Require().NoError(err)
		s.Require().EqualValues(model.URL{}, res)
	})

	s.Run("Remove urls by non-creator", func() {
		err := s.storage.RemoveUserURLs(s.ctx, foreignURLs)
		s.Require().NoError(err)

		expectedURL := s.fixtures.URLS[1].ToCanonical()
		res, err := s.storage.GetURL(s.ctx, 2)
		s.Require().NoError(err)
		s.Assert().EqualValues(expectedURL, res)
	})
}
