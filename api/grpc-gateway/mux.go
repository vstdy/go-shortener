package grpcgateway

import (
	"context"
	"fmt"

	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/encoding/protojson"

	"github.com/vstdy/go-shortener/pkg/grpc/url-service"
)

func NewMux(grpcSrvAddr string) (*runtime.ServeMux, error) {
	conn, err := grpc.Dial(
		grpcSrvAddr,
		grpc.WithUnaryInterceptor(
			grpc_opentracing.UnaryClientInterceptor(),
		),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("creating client connection failed: %v", err)
	}

	jsonPb := &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames:   true,
			EmitUnpopulated: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	}
	mux := runtime.NewServeMux(
		runtime.WithMetadata(metadataAnnotator),
		runtime.WithForwardResponseOption(httpResponseModifier),
		runtime.WithMarshalerOption(runtime.MIMEWildcard, jsonPb),
	)

	if err = urlService.RegisterURLServiceHandler(context.Background(), mux, conn); err != nil {
		return nil, fmt.Errorf("register handler failed: %v", err)
	}

	return mux, nil
}
