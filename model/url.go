package model

import "github.com/google/uuid"

type URL struct {
	ID     int       `json:"id"`
	UserID uuid.UUID `json:"user_id"`
	URL    string    `json:"url"`
}
