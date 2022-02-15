package model

import (
	"strconv"

	"github.com/google/uuid"

	"github.com/vstdy0/go-project/model"
)

type AddURLRequest struct {
	URL string `json:"url"`
}

// ToCanonical converts a API model to canonical model.
func (u AddURLRequest) ToCanonical(userID uuid.UUID) model.URL {
	obj := model.URL{
		UserID: userID,
		URL:    u.URL,
	}

	return obj
}

type AddURLResponse struct {
	Result string `json:"result"`
}

type (
	urlInBatch struct {
		CorrelationID string `json:"correlation_id"`
		OriginalURL   string `json:"original_url"`
	}

	AddURLsBatchRequest []urlInBatch
)

// ToCanonical converts a API model to canonical model.
func (u AddURLsBatchRequest) ToCanonical(userID uuid.UUID) ([]model.URL, error) {
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

type AddURLsBatchResponse struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

// NewURLsBatchFromCanonical creates a new URLS DB object from canonical model.
func NewURLsBatchFromCanonical(objs []model.URL, baseURL string) []AddURLsBatchResponse {
	var urls []AddURLsBatchResponse
	for _, url := range objs {
		urls = append(urls, AddURLsBatchResponse{
			CorrelationID: url.CorrelationID,
			ShortURL:      baseURL + "/" + strconv.Itoa(url.ID),
		})
	}

	return urls
}
