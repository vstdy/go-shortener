package model

import (
	"strconv"
	"strings"

	"github.com/google/uuid"

	"github.com/vstdy/go-shortener/model"
	"github.com/vstdy/go-shortener/pkg/grpc/url-service"
)

// newShortcut returns shortcut for object.
func newShortcut(objID int, baseURL string) string {
	return baseURL + "/gw/" + strconv.Itoa(objID)
}

// ShortenURL
// NewShortenURLReq creates new ShortenURLReq model from url.
func NewShortenURLReq(url string) *urlService.ShortenURLReq {
	return &urlService.ShortenURLReq{Url: url}
}

// ShortenURLReqToCanon converts gRPC model to canonical model.
func ShortenURLReqToCanon(in *urlService.ShortenURLReq, userID uuid.UUID) model.URL {
	return model.URL{
		UserID: userID,
		URL:    in.Url,
	}
}

// ShortenURLRespFromCanon converts canonical model to gRPC model.
func ShortenURLRespFromCanon(obj model.URL, baseURL string) *urlService.ShortenURLResp {
	return &urlService.ShortenURLResp{Result: newShortcut(obj.ID, baseURL)}
}

// ShortenURLsBatch
// NewShortenURLsBatchReq creates new ShortenURLsBatchReq model from slice of urls.
func NewShortenURLsBatchReq(urls []string) *urlService.ShortenURLsBatchReq {
	var units []*urlService.ShortenURLsBatchReq_UrlUnit
	for idx, url := range urls {
		units = append(units, &urlService.ShortenURLsBatchReq_UrlUnit{
			CorrelationId: strconv.Itoa(idx + 1),
			OriginalUrl:   url,
		})
	}

	return &urlService.ShortenURLsBatchReq{Request: units}
}

// ShortenURLsBatchReqToCanon converts gRPC model to canonical model.
func ShortenURLsBatchReqToCanon(in *urlService.ShortenURLsBatchReq, userID uuid.UUID) []model.URL {
	var objs []model.URL
	for _, unit := range in.GetRequest() {
		objs = append(objs, model.URL{
			CorrelationID: unit.GetCorrelationId(),
			UserID:        userID,
			URL:           unit.GetOriginalUrl(),
		})
	}

	return objs
}

// ShortenURLsBatchRespFromCanon converts canonical model to gRPC model.
func ShortenURLsBatchRespFromCanon(objs []model.URL, baseURL string) *urlService.ShortenURLsBatchResp {
	var resp []*urlService.ShortenURLsBatchResp_UrlUnit
	for _, obj := range objs {
		resp = append(resp, &urlService.ShortenURLsBatchResp_UrlUnit{
			CorrelationId: obj.CorrelationID,
			ShortUrl:      newShortcut(obj.ID, baseURL),
		})
	}

	return &urlService.ShortenURLsBatchResp{Response: resp}
}

// GetOriginalURL
// NewGetOrigURLReq creates new GetOrigURLReq model from url.
func NewGetOrigURLReq(url, baseURL string) *urlService.GetOrigURLReq {
	id := strings.TrimPrefix(url, baseURL+"/gw/")

	return &urlService.GetOrigURLReq{Id: id}
}

// GetUsersURLs
// GetUsersURLsRespFromCanon converts canonical model to gRPC model.
func GetUsersURLsRespFromCanon(objs []model.URL, baseURL string) *urlService.GetUsersURLsResp {
	var urls []*urlService.GetUsersURLsResp_UrlUnit
	for _, obj := range objs {
		urls = append(urls, &urlService.GetUsersURLsResp_UrlUnit{
			OriginalUrl: obj.URL,
			ShortUrl:    newShortcut(obj.ID, baseURL),
		})
	}

	return &urlService.GetUsersURLsResp{Response: urls}
}

// DeleteUserURLs
// NewDelUserURLsReq creates gRPC model to canonical model.
func NewDelUserURLsReq(ids []string) *urlService.DelUserURLsReq {
	return &urlService.DelUserURLsReq{Ids: ids}
}

// DelUserURLsReqToCanon converts gRPC model to canonical model.
func DelUserURLsReqToCanon(in *urlService.DelUserURLsReq, userID uuid.UUID) ([]model.URL, error) {
	var objs []model.URL
	for _, idRaw := range in.GetIds() {
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
