package schema

import (
	"github.com/google/uuid"

	"github.com/vstdy0/go-project/model"
)

type (
	URL struct {
		ID     int       `json:"id"`
		UserID uuid.UUID `json:"user_id"`
		URL    string    `json:"url"`
	}

	URLS []URL
)

// NewURLFromCanonical creates a new URL DB object from canonical model.
func NewURLFromCanonical(obj model.URL) URL {
	return URL{
		ID:     obj.ID,
		UserID: obj.UserID,
		URL:    obj.URL,
	}
}

// ToCanonical converts a DB object to canonical model.
func (u URL) ToCanonical() model.URL {
	obj := model.URL{
		ID:     u.ID,
		UserID: u.UserID,
		URL:    u.URL,
	}

	return obj
}

// ToCanonical converts a DB object to canonical model.
func (u URLS) ToCanonical() []model.URL {
	objs := make([]model.URL, 0, len(u))
	for _, dbObj := range u {
		obj := dbObj.ToCanonical()
		objs = append(objs, obj)
	}

	return objs
}
