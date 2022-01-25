package api

import (
	"encoding/json"
)

func (h Handler) jsonResponse(userID string, body []byte) ([]byte, error) {
	var urlRequest shortenURLRequest
	if err := json.Unmarshal(body, &urlRequest); err != nil {
		return nil, err
	}
	id, err := h.service.AddURL(userID, urlRequest.URL)
	if err != nil {
		return nil, err
	}
	res, err := json.Marshal(shortenURLResponse{Result: h.cfg.BaseURL + "/" + id})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (h Handler) plainResponse(userID string, body []byte) ([]byte, error) {
	id, err := h.service.AddURL(userID, string(body))
	if err != nil {
		return nil, err
	}

	return []byte(h.cfg.BaseURL + "/" + id), nil
}
