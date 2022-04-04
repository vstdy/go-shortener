package rest

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"

	"github.com/google/uuid"

	"github.com/vstdy0/go-shortener/api/rest/model"
	"github.com/vstdy0/go-shortener/pkg"
)

func (h Handler) jsonURLResponse(ctx context.Context, userID uuid.UUID, body []byte) ([]byte, error) {
	var req model.AddURLRequest

	if err := json.Unmarshal(body, &req); err != nil {
		return nil, err
	}

	obj := req.ToCanonical(userID)
	svcErr := h.service.AddURL(ctx, &obj)
	if svcErr != nil && !errors.Is(svcErr, pkg.ErrAlreadyExists) {
		return nil, svcErr
	}

	res, err := json.Marshal(model.NewURLRespFromCanon(obj, h.config.BaseURL))
	if err != nil {
		return nil, err
	}

	return res, svcErr
}

func (h Handler) plainURLResponse(ctx context.Context, userID uuid.UUID, body []byte) ([]byte, error) {
	req := model.NewURLReqFromStr(string(body))

	obj := req.ToCanonical(userID)
	svcErr := h.service.AddURL(ctx, &obj)
	if svcErr != nil && !errors.Is(svcErr, pkg.ErrAlreadyExists) {
		return nil, svcErr
	}

	return []byte(h.config.BaseURL + "/" + strconv.Itoa(obj.ID)), svcErr
}

func (h Handler) urlsBatchResponse(ctx context.Context, userID uuid.UUID, body []byte) ([]byte, error) {
	var batchReq model.AddURLsBatchReq
	if err := json.Unmarshal(body, &batchReq); err != nil {
		return nil, err
	}

	objs, err := batchReq.ToCanonical(userID)
	if err != nil {
		return nil, err
	}

	svcErr := h.service.AddURLsBatch(ctx, &objs)
	if svcErr != nil && !errors.Is(svcErr, pkg.ErrAlreadyExists) {
		return nil, svcErr
	}

	batchRes := model.NewURLsBatchRespFromCanon(objs, h.config.BaseURL)

	res, err := json.Marshal(batchRes)
	if err != nil {
		return nil, err
	}

	return res, svcErr
}
