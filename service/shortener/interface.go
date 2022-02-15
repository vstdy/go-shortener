package shortener

import (
	"context"
	"io"

	"github.com/google/uuid"

	"github.com/vstdy0/go-project/model"
)

type URLService interface {
	io.Closer

	AddURL(ctx context.Context, obj *model.URL) error
	AddBatchURLs(ctx context.Context, objs *[]model.URL) error
	GetURL(ctx context.Context, urlID int) (string, error)
	GetUserURLs(ctx context.Context, userID uuid.UUID) ([]model.URL, error)
	DeleteUserURLs(objs []model.URL) error
	Ping() error
}
