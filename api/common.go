package api

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/google/uuid"
)

func (h Handler) jsonResponse(ctx context.Context, userID uuid.UUID, body []byte) ([]byte, error) {
	var urlRequest shortenURLRequest
	if err := json.Unmarshal(body, &urlRequest); err != nil {
		return nil, err
	}
	id, err := h.service.AddURL(ctx, userID, urlRequest.URL)
	if err != nil {
		return nil, err
	}
	res, err := json.Marshal(shortenURLResponse{Result: h.cfg.BaseURL + "/" + strconv.Itoa(id)})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (h Handler) plainResponse(ctx context.Context, userID uuid.UUID, body []byte) ([]byte, error) {
	id, err := h.service.AddURL(ctx, userID, string(body))
	if err != nil {
		return nil, err
	}

	return []byte(h.cfg.BaseURL + "/" + strconv.Itoa(id)), nil
}
