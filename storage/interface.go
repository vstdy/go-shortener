package storage

import (
	"context"
	"io"

	"github.com/google/uuid"

	"github.com/vstdy0/go-project/model"
)

type URLStorage interface {
	io.Closer

	HasURL(ctx context.Context, urlID int) (bool, error)
	AddURLS(ctx context.Context, urls []model.URL) ([]model.URL, error)
	GetURL(ctx context.Context, urlID int) (model.URL, error)
	GetUserURLs(ctx context.Context, userID uuid.UUID) ([]model.URL, error)
	Ping() error
}
