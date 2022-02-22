package model

import "github.com/google/uuid"

// URL keeps url data.
type URL struct {
	ID            int
	CorrelationID string
	UserID        uuid.UUID
	URL           string
}
