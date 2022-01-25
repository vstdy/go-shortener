package storage

import (
	"github.com/vstdy0/go-project/model"
)

type URLStorage interface {
	Has(urlID string) bool
	Set(urlID, userID, url string) (string, error)
	Get(urlID string) string
	GetUserURLs(userID string) []model.URL
}
