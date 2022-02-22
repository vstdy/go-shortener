package shortener

import (
	"errors"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"

	"github.com/vstdy0/go-project/model"
	"github.com/vstdy0/go-project/pkg"
	storagemock "github.com/vstdy0/go-project/storage/mock"
)

func (s *TestSuite) TestURLService_AddURL() {
	type testCase struct {
		name         string
		prepareMocks func(URLStorageMock *storagemock.MockURLStorage) model.URL
		errExpected  bool
		errTarget    error
		errContains  string
	}

	testCases := []testCase{
		{
			name: "Fail: invalid input (empty url)",
			prepareMocks: func(URLStorageMock *storagemock.MockURLStorage) model.URL {
				return model.URL{
					UserID: uuid.New(),
					URL:    "",
				}
			},
			errExpected: true,
			errTarget:   pkg.ErrInvalidInput,
			errContains: "url",
		},
		{
			name: "Fail: invalid input (invalid url)",
			prepareMocks: func(URLStorageMock *storagemock.MockURLStorage) model.URL {
				return model.URL{
					UserID: uuid.New(),
					URL:    "htp//invalid-url.com/",
				}
			},
			errExpected: true,
			errTarget:   pkg.ErrInvalidInput,
			errContains: "url",
		},
		{
			name: "OK",
			prepareMocks: func(URLStorageMock *storagemock.MockURLStorage) model.URL {
				input := model.URL{
					UserID: uuid.New(),
					URL:    "https://extremely-lengthy-url.com/",
				}

				urls := []model.URL{
					{
						ID:     1,
						UserID: input.UserID,
						URL:    input.URL,
					},
				}

				URLStorageMock.EXPECT().
					AddURLs(gomock.Any(), []model.URL{input}).
					Return(urls, nil)

				return input
			},
			errExpected: false,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			input := tc.prepareMocks(s.stMock)

			err := s.svc.AddURL(s.ctx, &input)
			if tc.errExpected {
				s.Assert().Error(err)
				if tc.errTarget != nil {
					s.Assert().True(errors.Is(err, tc.errTarget))
				}
				if tc.errContains != "" {
					s.Assert().Contains(err.Error(), tc.errContains)
				}
				return
			}

			s.Assert().NoError(err)
		})
	}
}

func (s *TestSuite) TestURLService_AddBatchURLs() {
	type testCase struct {
		name         string
		prepareMocks func(URLStorageMock *storagemock.MockURLStorage) []model.URL
		errExpected  bool
		errTarget    error
		errContains  string
	}

	testCases := []testCase{
		{
			name: "Fail: invalid input (empty correlation_id)",
			prepareMocks: func(URLStorageMock *storagemock.MockURLStorage) []model.URL {
				return []model.URL{
					{
						CorrelationID: uuid.NewString(),
						UserID:        uuid.New(),
						URL:           "https://not-empty-valid-url.com/",
					},
					{
						CorrelationID: "",
						UserID:        uuid.New(),
						URL:           "https://extremely-lengthy-url.com/",
					},
				}
			},
			errExpected: true,
			errTarget:   pkg.ErrInvalidInput,
			errContains: "correlation_id",
		},
		{
			name: "Fail: invalid input (empty url)",
			prepareMocks: func(URLStorageMock *storagemock.MockURLStorage) []model.URL {
				return []model.URL{
					{
						CorrelationID: uuid.NewString(),
						UserID:        uuid.New(),
						URL:           "https://not-empty-valid-url.com/",
					},
					{
						CorrelationID: uuid.NewString(),
						UserID:        uuid.New(),
						URL:           "",
					},
				}
			},
			errExpected: true,
			errTarget:   pkg.ErrInvalidInput,
			errContains: "url",
		},
		{
			name: "Fail: invalid input (invalid url)",
			prepareMocks: func(URLStorageMock *storagemock.MockURLStorage) []model.URL {
				return []model.URL{
					{
						CorrelationID: uuid.NewString(),
						UserID:        uuid.New(),
						URL:           "https://not-empty-valid-url.com/",
					},
					{
						CorrelationID: uuid.NewString(),
						UserID:        uuid.New(),
						URL:           "htp//invalid-url.com/",
					},
				}
			},
			errExpected: true,
			errTarget:   pkg.ErrInvalidInput,
			errContains: "url",
		},
		{
			name: "OK",
			prepareMocks: func(URLStorageMock *storagemock.MockURLStorage) []model.URL {
				input := []model.URL{
					{
						CorrelationID: uuid.NewString(),
						UserID:        uuid.New(),
						URL:           "https://not-empty-valid-url.com/",
					},
					{
						CorrelationID: uuid.NewString(),
						UserID:        uuid.New(),
						URL:           "https://extremely-lengthy-url.com/",
					},
				}

				urls := []model.URL{
					{
						ID:            1,
						CorrelationID: input[0].CorrelationID,
						UserID:        input[0].UserID,
						URL:           input[0].URL,
					},
					{
						ID:            2,
						CorrelationID: input[1].CorrelationID,
						UserID:        input[1].UserID,
						URL:           input[1].URL,
					},
				}

				URLStorageMock.EXPECT().
					AddURLs(gomock.Any(), input).
					Return(urls, nil)

				return input
			},
			errExpected: false,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			input := tc.prepareMocks(s.stMock)

			err := s.svc.AddBatchURLs(s.ctx, &input)
			if tc.errExpected {
				s.Assert().Error(err)
				if tc.errTarget != nil {
					s.Assert().True(errors.Is(err, tc.errTarget))
				}
				if tc.errContains != "" {
					s.Assert().Contains(err.Error(), tc.errContains)
				}
				return
			}

			s.Assert().NoError(err)
		})
	}
}

func (s *TestSuite) TestURLService_GetURL() {
	type testCase struct {
		name         string
		prepareMocks func(URLStorageMock *storagemock.MockURLStorage) int
		errExpected  bool
		errTarget    error
		errContains  string
	}

	testCases := []testCase{
		{
			name: "Fail: invalid input (zero id)",
			prepareMocks: func(URLStorageMock *storagemock.MockURLStorage) int {
				return 0
			},
			errExpected: true,
			errTarget:   pkg.ErrInvalidInput,
			errContains: "id",
		},
		{
			name: "Fail: invalid input (negative id)",
			prepareMocks: func(URLStorageMock *storagemock.MockURLStorage) int {
				return -1
			},
			errExpected: true,
			errTarget:   pkg.ErrInvalidInput,
			errContains: "id",
		},
		{
			name: "OK",
			prepareMocks: func(URLStorageMock *storagemock.MockURLStorage) int {
				input := 1

				url := model.URL{
					ID:     1,
					UserID: uuid.New(),
					URL:    "https://extremely-lengthy-url.com/",
				}

				URLStorageMock.EXPECT().
					GetURL(gomock.Any(), input).
					Return(url, nil)

				return input
			},
			errExpected: false,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			input := tc.prepareMocks(s.stMock)

			_, err := s.svc.GetURL(s.ctx, input)
			if tc.errExpected {
				s.Assert().Error(err)
				if tc.errTarget != nil {
					s.Assert().True(errors.Is(err, tc.errTarget))
				}
				if tc.errContains != "" {
					s.Assert().Contains(err.Error(), tc.errContains)
				}
				return
			}

			s.Assert().NoError(err)
		})
	}
}

func (s *TestSuite) TestURLService_GetUserURL() {
	type testCase struct {
		name         string
		prepareMocks func(URLStorageMock *storagemock.MockURLStorage) uuid.UUID
		errExpected  bool
		errTarget    error
		errContains  string
	}

	testCases := []testCase{
		{
			name: "OK",
			prepareMocks: func(URLStorageMock *storagemock.MockURLStorage) uuid.UUID {
				input := uuid.New()

				urls := []model.URL{
					{
						ID:     1,
						UserID: input,
						URL:    "https://extremely-lengthy-url.com/",
					},
					{
						ID:     2,
						UserID: input,
						URL:    "https://another-lengthy-url.com/",
					},
				}

				URLStorageMock.EXPECT().
					GetUserURLs(gomock.Any(), input).
					Return(urls, nil)

				return input
			},
			errExpected: false,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			input := tc.prepareMocks(s.stMock)

			_, err := s.svc.GetUserURLs(s.ctx, input)
			if tc.errExpected {
				s.Assert().Error(err)
				if tc.errTarget != nil {
					s.Assert().True(errors.Is(err, tc.errTarget))
				}
				if tc.errContains != "" {
					s.Assert().Contains(err.Error(), tc.errContains)
				}
				return
			}

			s.Assert().NoError(err)
		})
	}
}

func (s *TestSuite) TestURLService_RemoveUserURLs() {
	type testCase struct {
		name         string
		prepareMocks func(URLStorageMock *storagemock.MockURLStorage) []model.URL
		errExpected  bool
		errTarget    error
		errContains  string
	}

	testCases := []testCase{
		{
			name: "Fail: invalid input (empty ids list)",
			prepareMocks: func(URLStorageMock *storagemock.MockURLStorage) []model.URL {
				return nil
			},
			errExpected: true,
			errTarget:   pkg.ErrInvalidInput,
			errContains: "ids",
		},
		{
			name: "OK with full buffer",
			prepareMocks: func(URLStorageMock *storagemock.MockURLStorage) []model.URL {
				input := []model.URL{
					{
						ID:     1,
						UserID: uuid.New(),
					},
					{
						ID:     2,
						UserID: uuid.New(),
					},
				}

				s.Add(1)

				URLStorageMock.EXPECT().
					RemoveUserURLs(gomock.Any(), input).
					Do(func(arg0, arg1 interface{}) { s.Done() }).
					Return(nil)

				return input
			},
			errExpected: false,
		},
		{
			name: "OK with buffer wipe timeout",
			prepareMocks: func(URLStorageMock *storagemock.MockURLStorage) []model.URL {
				input := []model.URL{
					{
						ID:     1,
						UserID: uuid.New(),
					},
				}

				s.Add(1)

				URLStorageMock.EXPECT().
					RemoveUserURLs(gomock.Any(), input).
					Do(func(arg0, arg1 interface{}) { s.Done() }).
					Return(nil)

				return input
			},
			errExpected: false,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			input := tc.prepareMocks(s.stMock)

			err := s.svc.RemoveUserURLs(input)
			if tc.errExpected {
				s.Assert().Error(err)
				if tc.errTarget != nil {
					s.Assert().True(errors.Is(err, tc.errTarget))
				}
				if tc.errContains != "" {
					s.Assert().Contains(err.Error(), tc.errContains)
				}
				return
			}

			s.Wait()

			s.Assert().NoError(err)
		})
	}
}
