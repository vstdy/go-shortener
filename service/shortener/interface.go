package shortener

import "github.com/vstdy0/go-project/model"

type URLService interface {
	AddURL(userID, id string) (string, error)
	GetURL(id string) string
	GetUserURLs(userID string) []model.URL
	GetUserID() int
}
