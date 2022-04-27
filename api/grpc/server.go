package grpc

import (
	"fmt"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	grpc_validator "github.com/grpc-ecosystem/go-grpc-middleware/validator"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"

	"github.com/vstdy/go-shortener/pkg/grpc/url-service"
	"github.com/vstdy/go-shortener/service/shortener"
)

var _ urlService.URLServiceServer = (*gRPCServer)(nil)

// gRPCServer implement gRPC services.
type gRPCServer struct {
	urlService.UnimplementedURLServiceServer

	config  Config
	service shortener.Service
}

// registerServices registers grpc services
func (srv *gRPCServer) registerServices(gRPCServer *grpc.Server) {
	urlService.RegisterURLServiceServer(gRPCServer, srv)
}

// Option specifies functional argument used by NewServer function.
type Option func(server *gRPCServer) error

// WithService sets service for server.
func WithService(service shortener.Service) Option {
	return func(srv *gRPCServer) error {
		srv.service = service

		return nil
	}
}

// WithConfig sets Config for server.
func WithConfig(config Config) Option {
	return func(srv *gRPCServer) error {
		srv.config = config

		return nil
	}
}

// NewServer creates a new configured gRPCServer object.
func NewServer(options ...Option) (*grpc.Server, error) {
	srv := &gRPCServer{
		config: NewDefaultConfig(),
	}
	for _, option := range options {
		if err := option(srv); err != nil {
			return nil, err
		}
	}

	if err := srv.config.Validate(); err != nil {
		return nil, fmt.Errorf("config validation: %w", err)
	}

	gRPCServer := grpc.NewServer(
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle: srv.config.KapMaxConnectionIdle,
			MaxConnectionAge:  srv.config.KapMaxConnectionAge,
			Time:              srv.config.KapTime,
			Timeout:           srv.config.KapTimeout,
		}),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_validator.UnaryServerInterceptor(),
			grpc_ctxtags.UnaryServerInterceptor(),
			grpc_opentracing.UnaryServerInterceptor(),
			grpc_recovery.UnaryServerInterceptor(),
			grpc_auth.UnaryServerInterceptor(metadataAuth(srv.config.SecretKey)),
			timeoutInterceptor(srv.config.Timeout),
		)),
	)
	srv.registerServices(gRPCServer)

	return gRPCServer, nil
}
