package shortener

import (
	"context"
	"errors"
	"github.com/vstdy0/go-project/pkg"

	"github.com/google/uuid"

	"github.com/vstdy0/go-project/model"
)

func (svc *Service) AddURL(ctx context.Context, url *model.URL) error {
	urlsModel, err := svc.storage.AddURLS(ctx, []model.URL{*url})
	if err != nil {
		if errors.Is(err, pkg.ErrIntegrityViolation) {
			url.ID = urlsModel[0].ID
		}
		return err
	}

	url.ID = urlsModel[0].ID

	return nil
}

func (svc *Service) AddBatchURLs(ctx context.Context, urls *[]model.URL) error {
	urlsModel, err := svc.storage.AddURLS(ctx, *urls)
	if err != nil {
		if errors.Is(err, pkg.ErrIntegrityViolation) {
			*urls = urlsModel
		}
		return err
	}

	*urls = urlsModel

	return nil
}

func (svc *Service) GetURL(ctx context.Context, id int) (string, error) {
	urlModel, err := svc.storage.GetURL(ctx, id)
	if err != nil {
		return "", err
	}

	return urlModel.URL, nil
}

func (svc *Service) GetUserURLs(ctx context.Context, userID uuid.UUID) ([]model.URL, error) {
	urls, err := svc.storage.GetUserURLs(ctx, userID)
	if err != nil {
		return nil, err
	}

	return urls, nil
}

func (svc *Service) Ping() error {
	if err := svc.storage.Ping(); err != nil {
		return err
	}

	return nil
}
