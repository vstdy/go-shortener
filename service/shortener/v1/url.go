package shortener

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/vstdy0/go-project/model"
	"github.com/vstdy0/go-project/pkg"
	"github.com/vstdy0/go-project/service/shortener/v1/validator"
)

// AddURL adds given object to storage.
func (svc *Service) AddURL(ctx context.Context, obj *model.URL) error {
	if err := validator.ValidateURL(obj.URL); err != nil {
		return fmt.Errorf("%w: url: %v", pkg.ErrInvalidInput, err)
	}

	objs, err := svc.storage.AddURLs(ctx, []model.URL{*obj})
	if err != nil {
		if errors.Is(err, pkg.ErrIntegrityViolation) {
			obj.ID = objs[0].ID
		}
		return err
	}

	obj.ID = objs[0].ID

	return nil
}

// AddBatchURLs adds given batch of objects to storage.
func (svc *Service) AddBatchURLs(ctx context.Context, objs *[]model.URL) error {
	for _, obj := range *objs {
		if err := validator.ValidateURL(obj.URL); err != nil {
			return fmt.Errorf("%w: url: %v", pkg.ErrInvalidInput, err)
		}
		if obj.CorrelationID == "" {
			return fmt.Errorf("%w: correlation_id: empty", pkg.ErrInvalidInput)
		}
	}

	addedObjs, err := svc.storage.AddURLs(ctx, *objs)
	if err != nil {
		if errors.Is(err, pkg.ErrIntegrityViolation) {
			*objs = addedObjs
		}
		return err
	}

	*objs = addedObjs

	return nil
}

// GetURL gets object with given id.
func (svc *Service) GetURL(ctx context.Context, id int) (string, error) {
	if id < 1 {
		return "", fmt.Errorf("%w: id: less than 1", pkg.ErrInvalidInput)
	}

	urlModel, err := svc.storage.GetURL(ctx, id)
	if err != nil {
		return "", err
	}

	return urlModel.URL, nil
}

// GetUserURLs gets current user objects.
func (svc *Service) GetUserURLs(ctx context.Context, userID uuid.UUID) ([]model.URL, error) {
	urls, err := svc.storage.GetUserURLs(ctx, userID)
	if err != nil {
		return nil, err
	}

	return urls, nil
}

// RemoveUserURLs removes current user objects with given ids.
func (svc *Service) RemoveUserURLs(objs []model.URL) error {
	if objs == nil {
		return fmt.Errorf("%w: ids: empty", pkg.ErrInvalidInput)
	}

	go func() {
		for _, obj := range objs {
			svc.delChan <- obj
		}
	}()

	return nil
}

// Ping verifies a connection to the database is still alive.
func (svc *Service) Ping() error {
	if err := svc.storage.Ping(); err != nil {
		return err
	}

	return nil
}
