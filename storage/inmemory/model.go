package inmemory

import "github.com/vstdy0/go-project/model"

type (
	URL struct {
		ID     string `json:"id"`
		UserID string `json:"user_id"`
		URL    string `json:"url"`
	}

	URLS []URL
)

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
