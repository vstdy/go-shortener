package shortener

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/vstdy/go-shortener/model"
	"github.com/vstdy/go-shortener/pkg"
	"github.com/vstdy/go-shortener/pkg/tracing"
	"github.com/vstdy/go-shortener/service/shortener/v1/validator"
)

// AddURL adds given object to storage.
func (svc *Service) AddURL(ctx context.Context, obj *model.URL) (err error) {
	ctx, span := tracing.StartSpanFromCtx(ctx, "shortener AddURL")
	defer tracing.FinishSpan(span, err)

	if err = validator.ValidateURL(obj.URL); err != nil {
		return fmt.Errorf("%w: url: %v", pkg.ErrInvalidInput, err)
	}

	objs, err := svc.storage.AddURLs(ctx, []model.URL{*obj})
	if err != nil {
		if errors.Is(err, pkg.ErrAlreadyExists) {
			obj.ID = objs[0].ID
		}
		return fmt.Errorf("shortener: AddURL: %w", err)
	}

	obj.ID = objs[0].ID

	return nil
}

// AddURLsBatch adds given batch of objects to storage.
func (svc *Service) AddURLsBatch(ctx context.Context, objs *[]model.URL) (err error) {
	ctx, span := tracing.StartSpanFromCtx(ctx, "shortener AddURLsBatch")
	defer tracing.FinishSpan(span, err)

	if *objs == nil {
		return pkg.ErrInvalidInput
	}

	for _, obj := range *objs {
		if err = validator.ValidateURL(obj.URL); err != nil {
			return fmt.Errorf("shortener: %w: url: %v", pkg.ErrInvalidInput, err)
		}
		if obj.CorrelationID == "" {
			return fmt.Errorf("shortener: %w: correlation_id: empty", pkg.ErrInvalidInput)
		}
	}

	addedObjs, err := svc.storage.AddURLs(ctx, *objs)
	if err != nil {
		if errors.Is(err, pkg.ErrAlreadyExists) {
			*objs = addedObjs
		}
		return fmt.Errorf("shortener: AddURLsBatch: %w", err)
	}

	*objs = addedObjs

	return nil
}

// GetURL gets object with given id.
func (svc *Service) GetURL(ctx context.Context, id int) (url string, err error) {
	ctx, span := tracing.StartSpanFromCtx(ctx, "shortener GetURL")
	defer tracing.FinishSpan(span, err)

	if id < 1 {
		return "", fmt.Errorf("shortener: GetURL: %w: id: less than 1", pkg.ErrInvalidInput)
	}

	urlModel, err := svc.storage.GetURL(ctx, id)
	if err != nil {
		return "", fmt.Errorf("shortener: GetURL: %w", err)
	}

	return urlModel.URL, nil
}

// GetUsersURLs gets current user objects.
func (svc *Service) GetUsersURLs(ctx context.Context, userID uuid.UUID) (objs []model.URL, err error) {
	ctx, span := tracing.StartSpanFromCtx(ctx, "shortener GetUsersURLs")
	defer tracing.FinishSpan(span, err)

	objs, err = svc.storage.GetUsersURLs(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("shortener: GetUsersURLs: %w", err)
	}

	return objs, nil
}

// RemoveUsersURLs removes current user objects with given ids.
func (svc *Service) RemoveUsersURLs(ctx context.Context, objs []model.URL) (err error) {
	_, span := tracing.StartSpanFromCtx(ctx, "shortener RemoveUsersURLs")
	defer tracing.FinishSpan(span, err)

	if len(objs) == 0 {
		return fmt.Errorf("shortener: RemoveUsersURLs: %w: ids: empty", pkg.ErrInvalidInput)
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
