// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             (unknown)
// source: api/grpc/url_service.proto

package urlService

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// URLServiceClient is the client API for URLService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type URLServiceClient interface {
	ShortenURL(ctx context.Context, in *ShortenURLReq, opts ...grpc.CallOption) (*ShortenURLResp, error)
	ShortenURLsBatch(ctx context.Context, in *ShortenURLsBatchReq, opts ...grpc.CallOption) (*ShortenURLsBatchResp, error)
	GetOriginalURL(ctx context.Context, in *GetOrigURLReq, opts ...grpc.CallOption) (*emptypb.Empty, error)
	GetUsersURLs(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*GetUsersURLsResp, error)
	DeleteUserURLs(ctx context.Context, in *DelUserURLsReq, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

type uRLServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewURLServiceClient(cc grpc.ClientConnInterface) URLServiceClient {
	return &uRLServiceClient{cc}
}

func (c *uRLServiceClient) ShortenURL(ctx context.Context, in *ShortenURLReq, opts ...grpc.CallOption) (*ShortenURLResp, error) {
	out := new(ShortenURLResp)
	err := c.cc.Invoke(ctx, "/urlService.URLService/ShortenURL", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *uRLServiceClient) ShortenURLsBatch(ctx context.Context, in *ShortenURLsBatchReq, opts ...grpc.CallOption) (*ShortenURLsBatchResp, error) {
	out := new(ShortenURLsBatchResp)
	err := c.cc.Invoke(ctx, "/urlService.URLService/ShortenURLsBatch", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *uRLServiceClient) GetOriginalURL(ctx context.Context, in *GetOrigURLReq, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/urlService.URLService/GetOriginalURL", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *uRLServiceClient) GetUsersURLs(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*GetUsersURLsResp, error) {
	out := new(GetUsersURLsResp)
	err := c.cc.Invoke(ctx, "/urlService.URLService/GetUsersURLs", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *uRLServiceClient) DeleteUserURLs(ctx context.Context, in *DelUserURLsReq, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/urlService.URLService/DeleteUserURLs", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// URLServiceServer is the server API for URLService service.
// All implementations must embed UnimplementedURLServiceServer
// for forward compatibility
type URLServiceServer interface {
	ShortenURL(context.Context, *ShortenURLReq) (*ShortenURLResp, error)
	ShortenURLsBatch(context.Context, *ShortenURLsBatchReq) (*ShortenURLsBatchResp, error)
	GetOriginalURL(context.Context, *GetOrigURLReq) (*emptypb.Empty, error)
	GetUsersURLs(context.Context, *emptypb.Empty) (*GetUsersURLsResp, error)
	DeleteUserURLs(context.Context, *DelUserURLsReq) (*emptypb.Empty, error)
	mustEmbedUnimplementedURLServiceServer()
}

// UnimplementedURLServiceServer must be embedded to have forward compatible implementations.
type UnimplementedURLServiceServer struct {
}

func (UnimplementedURLServiceServer) ShortenURL(context.Context, *ShortenURLReq) (*ShortenURLResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ShortenURL not implemented")
}
func (UnimplementedURLServiceServer) ShortenURLsBatch(context.Context, *ShortenURLsBatchReq) (*ShortenURLsBatchResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ShortenURLsBatch not implemented")
}
func (UnimplementedURLServiceServer) GetOriginalURL(context.Context, *GetOrigURLReq) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetOriginalURL not implemented")
}
func (UnimplementedURLServiceServer) GetUsersURLs(context.Context, *emptypb.Empty) (*GetUsersURLsResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUsersURLs not implemented")
}
func (UnimplementedURLServiceServer) DeleteUserURLs(context.Context, *DelUserURLsReq) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteUserURLs not implemented")
}
func (UnimplementedURLServiceServer) mustEmbedUnimplementedURLServiceServer() {}

// UnsafeURLServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to URLServiceServer will
// result in compilation errors.
type UnsafeURLServiceServer interface {
	mustEmbedUnimplementedURLServiceServer()
}

func RegisterURLServiceServer(s grpc.ServiceRegistrar, srv URLServiceServer) {
	s.RegisterService(&URLService_ServiceDesc, srv)
}

func _URLService_ShortenURL_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ShortenURLReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(URLServiceServer).ShortenURL(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/urlService.URLService/ShortenURL",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(URLServiceServer).ShortenURL(ctx, req.(*ShortenURLReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _URLService_ShortenURLsBatch_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ShortenURLsBatchReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(URLServiceServer).ShortenURLsBatch(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/urlService.URLService/ShortenURLsBatch",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(URLServiceServer).ShortenURLsBatch(ctx, req.(*ShortenURLsBatchReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _URLService_GetOriginalURL_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetOrigURLReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(URLServiceServer).GetOriginalURL(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/urlService.URLService/GetOriginalURL",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(URLServiceServer).GetOriginalURL(ctx, req.(*GetOrigURLReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _URLService_GetUsersURLs_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(URLServiceServer).GetUsersURLs(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/urlService.URLService/GetUsersURLs",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(URLServiceServer).GetUsersURLs(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _URLService_DeleteUserURLs_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DelUserURLsReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(URLServiceServer).DeleteUserURLs(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/urlService.URLService/DeleteUserURLs",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(URLServiceServer).DeleteUserURLs(ctx, req.(*DelUserURLsReq))
	}
	return interceptor(ctx, in, info, handler)
}

// URLService_ServiceDesc is the grpc.ServiceDesc for URLService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var URLService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "urlService.URLService",
	HandlerType: (*URLServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ShortenURL",
			Handler:    _URLService_ShortenURL_Handler,
		},
		{
			MethodName: "ShortenURLsBatch",
			Handler:    _URLService_ShortenURLsBatch_Handler,
		},
		{
			MethodName: "GetOriginalURL",
			Handler:    _URLService_GetOriginalURL_Handler,
		},
		{
			MethodName: "GetUsersURLs",
			Handler:    _URLService_GetUsersURLs_Handler,
		},
		{
			MethodName: "DeleteUserURLs",
			Handler:    _URLService_DeleteUserURLs_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/grpc/url_service.proto",
}
