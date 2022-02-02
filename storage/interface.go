package storage

import (
	"context"

	"github.com/google/uuid"

	"github.com/vstdy0/go-project/model"
)

type URLStorage interface {
	Has(ctx context.Context, urlID int) (bool, error)
	Set(ctx context.Context, urls []model.URL) ([]model.URL, error)
	Get(ctx context.Context, urlID int) (model.URL, error)
	GetUserURLs(ctx context.Context, userID uuid.UUID) ([]model.URL, error)
}
