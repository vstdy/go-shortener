package model

import (
	"strconv"

	"github.com/google/uuid"

	"github.com/vstdy0/go-project/model"
)

type URLRequest struct {
	URL string `json:"url"`
}

// ToCanonical converts a API model to canonical model.
func (u URLRequest) ToCanonical(userID uuid.UUID) model.URL {
	obj := model.URL{
		UserID: userID,
		URL:    u.URL,
	}

	return obj
}

type URLResponse struct {
	Result string `json:"result"`
}

type (
	urlInBatch struct {
		CorrelationID string `json:"correlation_id"`
		OriginalURL   string `json:"original_url"`
	}

	URLsBatchRequest []urlInBatch
)

// ToCanonical converts a API model to canonical model.
func (u URLsBatchRequest) ToCanonical(userID uuid.UUID) ([]model.URL, error) {
	var objs []model.URL
	for _, url := range u {
		objs = append(objs, model.URL{
			CorrelationID: url.CorrelationID,
			UserID:        userID,
			URL:           url.OriginalURL,
		})
	}

	return objs, nil
}

type URLsBatchResponse struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

// NewURLsBatchFromCanonical creates a new URLS DB object from canonical model.
func NewURLsBatchFromCanonical(objs []model.URL, baseURL string) []URLsBatchResponse {
	var urls []URLsBatchResponse
	for _, url := range objs {
		urls = append(urls, URLsBatchResponse{
			CorrelationID: url.CorrelationID,
			ShortURL:      baseURL + "/" + strconv.Itoa(url.ID),
		})
	}

	return urls
}
