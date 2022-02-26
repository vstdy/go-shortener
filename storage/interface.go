//go:generate mockgen -source=interface.go -destination=./mock/storage.go -package=storagemock
package storage

import (
	"context"
	"io"

	"github.com/google/uuid"

	"github.com/vstdy0/go-project/model"
)

type URLStorage interface {
	io.Closer

	// HasURL checks existence of the object with given id
	HasURL(ctx context.Context, urlID int) (bool, error)
	// AddURLs adds given objects to storage
	AddURLs(ctx context.Context, objs []model.URL) ([]model.URL, error)
	// GetURL gets object with given id
	GetURL(ctx context.Context, urlID int) (model.URL, error)
	// GetUserURLs gets current user objects
	GetUserURLs(ctx context.Context, userID uuid.UUID) ([]model.URL, error)
	// RemoveUserURLs removes current user objects with given ids
	RemoveUserURLs(ctx context.Context, objs []model.URL) error
	// Ping verifies a connection to the database is still alive.
	Ping() error
}
