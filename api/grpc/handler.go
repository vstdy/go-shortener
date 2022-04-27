package grpc

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/vstdy/go-shortener/api/grpc/model"
	"github.com/vstdy/go-shortener/pkg"
	"github.com/vstdy/go-shortener/pkg/grpc/url-service"
)

// ShortenURL creates shortcut for given url.
func (srv *gRPCServer) ShortenURL(
	ctx context.Context, in *urlService.ShortenURLReq) (
	*urlService.ShortenURLResp, error) {

	userID, ok := ctx.Value(userIDKey).(uuid.UUID)
	if !ok {
		return nil, status.Error(codes.Internal, "context: failed to retrieve user_id")
	}

	header := metadata.New(map[string]string{})

	obj := model.ShortenURLReqToCanon(in, userID)

	err := srv.service.AddURL(ctx, &obj)
	if err != nil {
		if errors.Is(err, pkg.ErrInvalidInput) {
			return nil, status.Error(codes.InvalidArgument, pkg.ErrInvalidInput.Error())
		}

		if !errors.Is(err, pkg.ErrAlreadyExists) {
			return nil, status.Error(codes.Internal, err.Error())
		}

		header.Set("x-http-code", strconv.Itoa(http.StatusConflict))
	}

	out := model.ShortenURLRespFromCanon(obj, srv.config.BaseURL)

	header.Append("x-http-code", strconv.Itoa(http.StatusCreated))
	grpc.SetHeader(ctx, header)

	return out, nil
}

// ShortenURLsBatch creates shortcuts for given urls batch.
func (srv *gRPCServer) ShortenURLsBatch(
	ctx context.Context, in *urlService.ShortenURLsBatchReq) (
	*urlService.ShortenURLsBatchResp, error) {

	header := metadata.New(map[string]string{})

	userID, ok := ctx.Value(userIDKey).(uuid.UUID)
	if !ok {
		return nil, status.Error(codes.Internal, "context: failed to retrieve user_id")
	}

	objs := model.ShortenURLsBatchReqToCanon(in, userID)

	err := srv.service.AddURLsBatch(ctx, &objs)
	if err != nil {
		if errors.Is(err, pkg.ErrInvalidInput) {
			return nil, status.Error(codes.InvalidArgument, pkg.ErrInvalidInput.Error())
		}

		if !errors.Is(err, pkg.ErrAlreadyExists) {
			return nil, status.Error(codes.Internal, err.Error())
		}

		header.Set("x-http-code", strconv.Itoa(http.StatusConflict))
	}

	out := model.ShortenURLsBatchRespFromCanon(objs, srv.config.BaseURL)

	header.Append("x-http-code", strconv.Itoa(http.StatusCreated))
	grpc.SetHeader(ctx, header)

	return out, nil
}

// GetOriginalURL returns origin url from shortcut.
func (srv *gRPCServer) GetOriginalURL(
	ctx context.Context, in *urlService.GetOrigURLReq) (
	*emptypb.Empty, error) {

	header := metadata.New(map[string]string{})

	urlID, err := strconv.Atoi(in.GetId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, pkg.ErrInvalidInput.Error())
	}

	url, err := srv.service.GetURL(ctx, urlID)
	if err != nil {
		if errors.Is(err, pkg.ErrInvalidInput) {
			return nil, status.Error(codes.InvalidArgument, pkg.ErrInvalidInput.Error())
		}

		return nil, status.Error(codes.Internal, err.Error())
	}

	if url == "" {
		header.Set("x-http-code", strconv.Itoa(http.StatusGone))
		grpc.SetHeader(ctx, header)
		return &emptypb.Empty{}, nil
	}

	header.Set("location", url)
	header.Set("x-http-code", strconv.Itoa(http.StatusTemporaryRedirect))
	grpc.SetHeader(ctx, header)

	return &emptypb.Empty{}, nil
}

// GetUsersURLs returns urls created by current user.
func (srv *gRPCServer) GetUsersURLs(
	ctx context.Context, _ *emptypb.Empty) (
	*urlService.GetUsersURLsResp, error) {

	userID, ok := ctx.Value(userIDKey).(uuid.UUID)
	if !ok {
		return nil, status.Error(codes.Internal, "context: failed to retrieve user_id")
	}

	urls, err := srv.service.GetUsersURLs(ctx, userID)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	out := model.GetUsersURLsRespFromCanon(urls, srv.config.BaseURL)

	return out, nil
}

// DeleteUserURLs removes urls created by current user.
func (srv *gRPCServer) DeleteUserURLs(
	ctx context.Context, in *urlService.DelUserURLsReq) (
	*emptypb.Empty, error) {

	userID, ok := ctx.Value(userIDKey).(uuid.UUID)
	if !ok {
		return nil, status.Error(codes.Internal, "context: failed to retrieve user_id")
	}

	objs, err := model.DelUserURLsReqToCanon(in, userID)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, pkg.ErrInvalidInput.Error())
	}

	err = srv.service.RemoveUsersURLs(ctx, objs)
	if err != nil {
		if errors.Is(err, pkg.ErrInvalidInput) {
			return nil, status.Error(codes.InvalidArgument, pkg.ErrInvalidInput.Error())
		}

		return nil, status.Error(codes.Internal, err.Error())
	}

	header := metadata.New(map[string]string{"x-http-code": strconv.Itoa(http.StatusAccepted)})
	grpc.SetHeader(ctx, header)

	return &emptypb.Empty{}, nil
}
