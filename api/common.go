package api

import (
	"encoding/json"
)

func (h Handler) jsonResponse(body []byte) ([]byte, error) {
	var urlRequest ShortenURLRequest
	if err := json.Unmarshal(body, &urlRequest); err != nil {
		return nil, err
	}
	id, err := h.service.AddURL(urlRequest.URL)
	if err != nil {
		return nil, err
	}
	res, err := json.Marshal(ShortenURLResponse{Result: h.cfg.BaseURL + "/" + id})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (h Handler) plainResponse(body []byte) ([]byte, error) {
	id, err := h.service.AddURL(string(body))
	if err != nil {
		return nil, err
	}

	return []byte(h.cfg.BaseURL + "/" + id), nil
}
