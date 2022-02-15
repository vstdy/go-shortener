package shortener

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"github.com/vstdy0/go-project/model"
	"github.com/vstdy0/go-project/pkg"
)

func (svc *Service) AddURL(ctx context.Context, obj *model.URL) error {
	objs, err := svc.storage.AddURLS(ctx, []model.URL{*obj})
	if err != nil {
		if errors.Is(err, pkg.ErrIntegrityViolation) {
			obj.ID = objs[0].ID
		}
		return err
	}

	obj.ID = objs[0].ID

	return nil
}

func (svc *Service) AddBatchURLs(ctx context.Context, objs *[]model.URL) error {
	addedObjs, err := svc.storage.AddURLS(ctx, *objs)
	if err != nil {
		if errors.Is(err, pkg.ErrIntegrityViolation) {
			*objs = addedObjs
		}
		return err
	}

	*objs = addedObjs

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

func (svc *Service) DeleteUserURLs(objs []model.URL) error {
	go func() {
		for _, obj := range objs {
			svc.delChan <- obj
		}
	}()

	return nil
}

func (svc *Service) Ping() error {
	if err := svc.storage.Ping(); err != nil {
		return err
	}

	return nil
}
