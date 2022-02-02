package shortener

import (
	"context"

	"github.com/google/uuid"

	"github.com/vstdy0/go-project/model"
)

type URLService interface {
	AddURL(ctx context.Context, url *model.URL) error
	AddBatchURLs(ctx context.Context, urls *[]model.URL) error
	GetURL(ctx context.Context, urlID int) (string, error)
	GetUserURLs(ctx context.Context, userID uuid.UUID) ([]model.URL, error)
}
