package model

import (
	"strconv"

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
