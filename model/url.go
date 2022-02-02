package model

import "github.com/google/uuid"

type URL struct {
	ID            int
	CorrelationID string
	UserID        uuid.UUID
	URL           string
}
