package shortener

import (
	"context"

	"github.com/google/uuid"

	"github.com/vstdy0/go-project/model"
)

type URLService interface {
	AddURL(ctx context.Context, userID uuid.UUID, url string) (int, error)
	GetURL(ctx context.Context, urlID int) (string, error)
	GetUserURLs(ctx context.Context, userID uuid.UUID) ([]model.URL, error)
}
