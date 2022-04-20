package model

import (
	"strconv"

	"github.com/google/uuid"

	"github.com/vstdy0/go-shortener/model"
)

// newShortcut returns shortcut for object.
func newShortcut(objID int, baseURL string) string {
	return baseURL + "/" + strconv.Itoa(objID)
}

type AddURLRequest struct {
	URL string `json:"url"`
}

// ToCanonical converts API model to canonical model.
func (u AddURLRequest) ToCanonical(userID uuid.UUID) model.URL {
	obj := model.URL{
		UserID: userID,
		URL:    u.URL,
	}

	return obj
}

// NewURLReqFromStr creates AddURLRequest object from string.
func NewURLReqFromStr(url string) AddURLRequest {
	return AddURLRequest{URL: url}
}

type AddURLResponse struct {
	Result string `json:"result"`
}

// NewURLRespFromCanon creates AddURLResponse object from canonical model.
func NewURLRespFromCanon(obj model.URL, baseURL string) AddURLResponse {
	return AddURLResponse{Result: newShortcut(obj.ID, baseURL)}
}

type (
	urlInBatch struct {
		CorrelationID string `json:"correlation_id"`
		OriginalURL   string `json:"original_url"`
	}

	AddURLsBatchReq []urlInBatch
)

// ToCanonical converts API model to canonical model.
func (u AddURLsBatchReq) ToCanonical(userID uuid.UUID) ([]model.URL, error) {
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

type AddURLsBatchResp struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

// NewURLsBatchRespFromCanon creates array of AddURLsBatchResp objects
// from array of canonical models.
func NewURLsBatchRespFromCanon(objs []model.URL, baseURL string) []AddURLsBatchResp {
	var urls []AddURLsBatchResp
	for _, obj := range objs {
		urls = append(urls, AddURLsBatchResp{
			CorrelationID: obj.CorrelationID,
			ShortURL:      newShortcut(obj.ID, baseURL),
		})
	}

	return urls
}

type UserURL struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

// NewUserURLsFromCanon creates array of UserURL objects
// from array of canonical models.
func NewUserURLsFromCanon(objs []model.URL, baseURL string) []UserURL {
	var userURLs []UserURL
	for _, obj := range objs {
		userURLs = append(userURLs,
			UserURL{
				ShortURL:    newShortcut(obj.ID, baseURL),
				OriginalURL: obj.URL,
			})
	}

	return userURLs
}

// URLsToDelToCanon creates array of canonical models from array of ids.
func URLsToDelToCanon(ids []string, userID uuid.UUID) ([]model.URL, error) {
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
