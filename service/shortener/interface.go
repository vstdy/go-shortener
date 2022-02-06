package shortener

import (
	"context"
	"io"

	"github.com/google/uuid"

	"github.com/vstdy0/go-project/model"
)

type URLService interface {
	io.Closer

	AddURL(ctx context.Context, url *model.URL) (error, error)
	AddBatchURLs(ctx context.Context, urls *[]model.URL) (error, error)
	GetURL(ctx context.Context, urlID int) (string, error)
	GetUserURLs(ctx context.Context, userID uuid.UUID) ([]model.URL, error)
	GetPing() error
}
