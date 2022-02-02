package api

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/google/uuid"

	"github.com/vstdy0/go-project/api/model"
)

func (h Handler) jsonURLResponse(ctx context.Context, userID uuid.UUID, body []byte) ([]byte, error) {
	var req model.URLRequest

	if err := json.Unmarshal(body, &req); err != nil {
		return nil, err
	}

	obj := req.ToCanonical(userID)
	if err := h.service.AddURL(ctx, &obj); err != nil {
		return nil, err
	}

	res, err := json.Marshal(model.URLResponse{Result: h.cfg.BaseURL + "/" + strconv.Itoa(obj.ID)})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (h Handler) plainURLResponse(ctx context.Context, userID uuid.UUID, body []byte) ([]byte, error) {
	req := model.URLRequest{
		URL: string(body),
	}

	obj := req.ToCanonical(userID)
	if err := h.service.AddURL(ctx, &obj); err != nil {
		return nil, err
	}

	return []byte(h.cfg.BaseURL + "/" + strconv.Itoa(obj.ID)), nil
}

func (h Handler) urlsBatchResponse(ctx context.Context, userID uuid.UUID, body []byte) ([]byte, error) {
	var batchReq model.URLsBatchRequest
	if err := json.Unmarshal(body, &batchReq); err != nil {
		return nil, err
	}

	objs, err := batchReq.ToCanonical(userID)
	if err != nil {
		return nil, err
	}

	if err = h.service.AddBatchURLs(ctx, &objs); err != nil {
		return nil, err
	}

	batchRes := model.NewURLsBatchFromCanonical(objs, h.cfg.BaseURL)

	res, err := json.Marshal(batchRes)
	if err != nil {
		return nil, err
	}

	return res, nil
}
