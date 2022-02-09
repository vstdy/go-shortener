package api

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"

	"github.com/google/uuid"

	"github.com/vstdy0/go-project/api/model"
	"github.com/vstdy0/go-project/pkg"
)

func (h Handler) jsonURLResponse(ctx context.Context, userID uuid.UUID, body []byte) ([]byte, error) {
	var req model.URLRequest

	if err := json.Unmarshal(body, &req); err != nil {
		return nil, err
	}

	obj := req.ToCanonical(userID)
	svcErr := h.service.AddURL(ctx, &obj)
	if svcErr != nil && !errors.Is(svcErr, pkg.ErrIntegrityViolation) {
		return nil, svcErr
	}

	res, err := json.Marshal(model.URLResponse{Result: h.cfg.BaseURL + "/" + strconv.Itoa(obj.ID)})
	if err != nil {
		return nil, err
	}

	return res, svcErr
}

func (h Handler) plainURLResponse(ctx context.Context, userID uuid.UUID, body []byte) ([]byte, error) {
	req := model.URLRequest{
		URL: string(body),
	}

	obj := req.ToCanonical(userID)
	svcErr := h.service.AddURL(ctx, &obj)
	if svcErr != nil && !errors.Is(svcErr, pkg.ErrIntegrityViolation) {
		return nil, svcErr
	}

	return []byte(h.cfg.BaseURL + "/" + strconv.Itoa(obj.ID)), svcErr
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

	svcErr := h.service.AddBatchURLs(ctx, &objs)
	if svcErr != nil && !errors.Is(svcErr, pkg.ErrIntegrityViolation) {
		return nil, svcErr
	}

	batchRes := model.NewURLsBatchFromCanonical(objs, h.cfg.BaseURL)

	res, err := json.Marshal(batchRes)
	if err != nil {
		return nil, err
	}

	return res, svcErr
}
