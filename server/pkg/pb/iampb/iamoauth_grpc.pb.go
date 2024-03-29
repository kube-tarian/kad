// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.12
// source: proto/oauthproto/oauth.proto

package oauthproto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// OauthServiceClient is the client API for OauthService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type OauthServiceClient interface {
	CreateOauthClient(ctx context.Context, in *OauthClientRequest, opts ...grpc.CallOption) (*OauthClientResponse, error)
	GetOauthToken(ctx context.Context, in *OauthTokenRequest, opts ...grpc.CallOption) (*OauthTokenResponse, error)
	ValidateOauthToken(ctx context.Context, in *ValidateOauthTokenRequest, opts ...grpc.CallOption) (*ValidateOauthTokenResponse, error)
	CreateClientCredentialsClient(ctx context.Context, in *CreateClientCredentialsClientRequest, opts ...grpc.CallOption) (*CreateClientCredentialsClientResponse, error)
}

type oauthServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewOauthServiceClient(cc grpc.ClientConnInterface) OauthServiceClient {
	return &oauthServiceClient{cc}
}

func (c *oauthServiceClient) CreateOauthClient(ctx context.Context, in *OauthClientRequest, opts ...grpc.CallOption) (*OauthClientResponse, error) {
	out := new(OauthClientResponse)
	err := c.cc.Invoke(ctx, "/OauthService.OauthService/CreateOauthClient", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *oauthServiceClient) GetOauthToken(ctx context.Context, in *OauthTokenRequest, opts ...grpc.CallOption) (*OauthTokenResponse, error) {
	out := new(OauthTokenResponse)
	err := c.cc.Invoke(ctx, "/OauthService.OauthService/GetOauthToken", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *oauthServiceClient) ValidateOauthToken(ctx context.Context, in *ValidateOauthTokenRequest, opts ...grpc.CallOption) (*ValidateOauthTokenResponse, error) {
	out := new(ValidateOauthTokenResponse)
	err := c.cc.Invoke(ctx, "/OauthService.OauthService/ValidateOauthToken", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *oauthServiceClient) CreateClientCredentialsClient(ctx context.Context, in *CreateClientCredentialsClientRequest, opts ...grpc.CallOption) (*CreateClientCredentialsClientResponse, error) {
	out := new(CreateClientCredentialsClientResponse)
	err := c.cc.Invoke(ctx, "/OauthService.OauthService/CreateClientCredentialsClient", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// OauthServiceServer is the server API for OauthService service.
// All implementations must embed UnimplementedOauthServiceServer
// for forward compatibility
type OauthServiceServer interface {
	CreateOauthClient(context.Context, *OauthClientRequest) (*OauthClientResponse, error)
	GetOauthToken(context.Context, *OauthTokenRequest) (*OauthTokenResponse, error)
	ValidateOauthToken(context.Context, *ValidateOauthTokenRequest) (*ValidateOauthTokenResponse, error)
	CreateClientCredentialsClient(context.Context, *CreateClientCredentialsClientRequest) (*CreateClientCredentialsClientResponse, error)
	mustEmbedUnimplementedOauthServiceServer()
}

// UnimplementedOauthServiceServer must be embedded to have forward compatible implementations.
type UnimplementedOauthServiceServer struct {
}

func (UnimplementedOauthServiceServer) CreateOauthClient(context.Context, *OauthClientRequest) (*OauthClientResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateOauthClient not implemented")
}
func (UnimplementedOauthServiceServer) GetOauthToken(context.Context, *OauthTokenRequest) (*OauthTokenResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetOauthToken not implemented")
}
func (UnimplementedOauthServiceServer) ValidateOauthToken(context.Context, *ValidateOauthTokenRequest) (*ValidateOauthTokenResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ValidateOauthToken not implemented")
}
func (UnimplementedOauthServiceServer) CreateClientCredentialsClient(context.Context, *CreateClientCredentialsClientRequest) (*CreateClientCredentialsClientResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateClientCredentialsClient not implemented")
}
func (UnimplementedOauthServiceServer) mustEmbedUnimplementedOauthServiceServer() {}

// UnsafeOauthServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to OauthServiceServer will
// result in compilation errors.
type UnsafeOauthServiceServer interface {
	mustEmbedUnimplementedOauthServiceServer()
}

func RegisterOauthServiceServer(s grpc.ServiceRegistrar, srv OauthServiceServer) {
	s.RegisterService(&OauthService_ServiceDesc, srv)
}

func _OauthService_CreateOauthClient_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(OauthClientRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OauthServiceServer).CreateOauthClient(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/OauthService.OauthService/CreateOauthClient",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OauthServiceServer).CreateOauthClient(ctx, req.(*OauthClientRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _OauthService_GetOauthToken_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(OauthTokenRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OauthServiceServer).GetOauthToken(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/OauthService.OauthService/GetOauthToken",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OauthServiceServer).GetOauthToken(ctx, req.(*OauthTokenRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _OauthService_ValidateOauthToken_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ValidateOauthTokenRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OauthServiceServer).ValidateOauthToken(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/OauthService.OauthService/ValidateOauthToken",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OauthServiceServer).ValidateOauthToken(ctx, req.(*ValidateOauthTokenRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _OauthService_CreateClientCredentialsClient_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateClientCredentialsClientRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OauthServiceServer).CreateClientCredentialsClient(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/OauthService.OauthService/CreateClientCredentialsClient",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OauthServiceServer).CreateClientCredentialsClient(ctx, req.(*CreateClientCredentialsClientRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// OauthService_ServiceDesc is the grpc.ServiceDesc for OauthService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var OauthService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "OauthService.OauthService",
	HandlerType: (*OauthServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateOauthClient",
			Handler:    _OauthService_CreateOauthClient_Handler,
		},
		{
			MethodName: "GetOauthToken",
			Handler:    _OauthService_GetOauthToken_Handler,
		},
		{
			MethodName: "ValidateOauthToken",
			Handler:    _OauthService_ValidateOauthToken_Handler,
		},
		{
			MethodName: "CreateClientCredentialsClient",
			Handler:    _OauthService_CreateClientCredentialsClient_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/oauthproto/oauth.proto",
}
