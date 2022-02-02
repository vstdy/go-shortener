package schema

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"

	"github.com/vstdy0/go-project/model"
)

type (
	URL struct {
		bun.BaseModel `bun:"url,alias:u"`
		ID            int       `bun:"id,pk,autoincrement"`
		CorrelationID string    `bun:"-"`
		UserID        uuid.UUID `bun:"user_id,type:uuid,notnull"`
		URL           string    `bun:"url,notnull"`
	}

	URLS []URL
)

// NewURLsFromCanonical creates new list of URL DB objects from canonical model.
func NewURLsFromCanonical(objs []model.URL) URLS {
	var urls URLS
	for _, url := range objs {
		urls = append(urls, URL{
			ID:            url.ID,
			CorrelationID: url.CorrelationID,
			UserID:        url.UserID,
			URL:           url.URL,
		})
	}

	return urls
}

// ToCanonical converts a DB object to canonical model.
func (u URL) ToCanonical() model.URL {
	obj := model.URL{
		ID:            u.ID,
		CorrelationID: u.CorrelationID,
		UserID:        u.UserID,
		URL:           u.URL,
	}

	return obj
}

// ToCanonical converts a DB object to canonical model.
func (u URLS) ToCanonical() ([]model.URL, error) {
	objs := make([]model.URL, 0, len(u))
	for _, dbObj := range u {
		obj := dbObj.ToCanonical()
		objs = append(objs, obj)
	}

	return objs, nil
}
