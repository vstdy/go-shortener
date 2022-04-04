//go:generate mockgen -source=interface.go -destination=./mock/service.go -package=servicemock
package shortener

import (
	"context"
	"io"

	"github.com/google/uuid"

	"github.com/vstdy0/go-shortener/model"
)

type Service interface {
	io.Closer

	// AddURL adds given object to storage.
	AddURL(ctx context.Context, obj *model.URL) error
	// AddURLsBatch adds given batch of objects to storage.
	AddURLsBatch(ctx context.Context, objs *[]model.URL) error
	// GetURL gets object with given id.
	GetURL(ctx context.Context, urlID int) (string, error)
	// GetUserURLs gets current user objects.
	GetUserURLs(ctx context.Context, userID uuid.UUID) ([]model.URL, error)
	// RemoveUserURLs removes current user objects with given ids.
	RemoveUserURLs(ctx context.Context, objs []model.URL) error
	// Ping verifies a connection to the database is still alive.
	Ping() error
}
