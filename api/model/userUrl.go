package model

import (
	"strconv"

	"github.com/google/uuid"

	"github.com/vstdy0/go-project/model"
)

type UserURL struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

func UserURLsFromCanonical(urls []model.URL, baseURL string) []UserURL {
	var userURLs []UserURL
	for _, v := range urls {
		userURLs = append(userURLs, UserURL{ShortURL: baseURL + "/" + strconv.Itoa(v.ID), OriginalURL: v.URL})
	}

	return userURLs
}

func URLsToDeleteToCanonical(ids []string, userID uuid.UUID) ([]model.URL, error) {
	var objs []model.URL
	for _, idRaw := range ids {
		id, err := strconv.Atoi(idRaw)
		if err != nil {
			return nil, err
		}
		objs = append(objs, model.URL{
			ID:     id,
			UserID: userID,
		})
	}

	return objs, nil
}
