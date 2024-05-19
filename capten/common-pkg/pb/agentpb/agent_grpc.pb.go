// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v5.26.1
// source: agent.proto

package agentpb

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

const (
	Agent_Ping_FullMethodName                   = "/agentpb.Agent/Ping"
	Agent_StoreCredential_FullMethodName        = "/agentpb.Agent/StoreCredential"
	Agent_ConfigureVaultSecret_FullMethodName   = "/agentpb.Agent/ConfigureVaultSecret"
	Agent_CreateVaultRole_FullMethodName        = "/agentpb.Agent/CreateVaultRole"
	Agent_UpdateVaultRole_FullMethodName        = "/agentpb.Agent/UpdateVaultRole"
	Agent_DeleteVaultRole_FullMethodName        = "/agentpb.Agent/DeleteVaultRole"
	Agent_SyncApp_FullMethodName                = "/agentpb.Agent/SyncApp"
	Agent_GetClusterApps_FullMethodName         = "/agentpb.Agent/GetClusterApps"
	Agent_GetClusterAppLaunches_FullMethodName  = "/agentpb.Agent/GetClusterAppLaunches"
	Agent_ConfigureAppSSO_FullMethodName        = "/agentpb.Agent/ConfigureAppSSO"
	Agent_GetClusterAppConfig_FullMethodName    = "/agentpb.Agent/GetClusterAppConfig"
	Agent_GetClusterAppValues_FullMethodName    = "/agentpb.Agent/GetClusterAppValues"
	Agent_GetClusterGlobalValues_FullMethodName = "/agentpb.Agent/GetClusterGlobalValues"
	Agent_DeployDefaultApps_FullMethodName      = "/agentpb.Agent/DeployDefaultApps"
	Agent_GetDefaultAppsStatus_FullMethodName   = "/agentpb.Agent/GetDefaultAppsStatus"
)

// AgentClient is the client API for Agent service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AgentClient interface {
	Ping(ctx context.Context, in *PingRequest, opts ...grpc.CallOption) (*PingResponse, error)
	StoreCredential(ctx context.Context, in *StoreCredentialRequest, opts ...grpc.CallOption) (*StoreCredentialResponse, error)
	ConfigureVaultSecret(ctx context.Context, in *ConfigureVaultSecretRequest, opts ...grpc.CallOption) (*ConfigureVaultSecretResponse, error)
	CreateVaultRole(ctx context.Context, in *CreateVaultRoleRequest, opts ...grpc.CallOption) (*CreateVaultRoleResponse, error)
	UpdateVaultRole(ctx context.Context, in *UpdateVaultRoleRequest, opts ...grpc.CallOption) (*UpdateVaultRoleResponse, error)
	DeleteVaultRole(ctx context.Context, in *DeleteVaultRoleRequest, opts ...grpc.CallOption) (*DeleteVaultRoleResponse, error)
	SyncApp(ctx context.Context, in *SyncAppRequest, opts ...grpc.CallOption) (*SyncAppResponse, error)
	GetClusterApps(ctx context.Context, in *GetClusterAppsRequest, opts ...grpc.CallOption) (*GetClusterAppsResponse, error)
	GetClusterAppLaunches(ctx context.Context, in *GetClusterAppLaunchesRequest, opts ...grpc.CallOption) (*GetClusterAppLaunchesResponse, error)
	ConfigureAppSSO(ctx context.Context, in *ConfigureAppSSORequest, opts ...grpc.CallOption) (*ConfigureAppSSOResponse, error)
	GetClusterAppConfig(ctx context.Context, in *GetClusterAppConfigRequest, opts ...grpc.CallOption) (*GetClusterAppConfigResponse, error)
	GetClusterAppValues(ctx context.Context, in *GetClusterAppValuesRequest, opts ...grpc.CallOption) (*GetClusterAppValuesResponse, error)
	GetClusterGlobalValues(ctx context.Context, in *GetClusterGlobalValuesRequest, opts ...grpc.CallOption) (*GetClusterGlobalValuesResponse, error)
	DeployDefaultApps(ctx context.Context, in *DeployDefaultAppsRequest, opts ...grpc.CallOption) (*DeployDefaultAppsResponse, error)
	GetDefaultAppsStatus(ctx context.Context, in *GetDefaultAppsStatusRequest, opts ...grpc.CallOption) (*GetDefaultAppsStatusResponse, error)
}

type agentClient struct {
	cc grpc.ClientConnInterface
}

func NewAgentClient(cc grpc.ClientConnInterface) AgentClient {
	return &agentClient{cc}
}

func (c *agentClient) Ping(ctx context.Context, in *PingRequest, opts ...grpc.CallOption) (*PingResponse, error) {
	out := new(PingResponse)
	err := c.cc.Invoke(ctx, Agent_Ping_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *agentClient) StoreCredential(ctx context.Context, in *StoreCredentialRequest, opts ...grpc.CallOption) (*StoreCredentialResponse, error) {
	out := new(StoreCredentialResponse)
	err := c.cc.Invoke(ctx, Agent_StoreCredential_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *agentClient) ConfigureVaultSecret(ctx context.Context, in *ConfigureVaultSecretRequest, opts ...grpc.CallOption) (*ConfigureVaultSecretResponse, error) {
	out := new(ConfigureVaultSecretResponse)
	err := c.cc.Invoke(ctx, Agent_ConfigureVaultSecret_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *agentClient) CreateVaultRole(ctx context.Context, in *CreateVaultRoleRequest, opts ...grpc.CallOption) (*CreateVaultRoleResponse, error) {
	out := new(CreateVaultRoleResponse)
	err := c.cc.Invoke(ctx, Agent_CreateVaultRole_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *agentClient) UpdateVaultRole(ctx context.Context, in *UpdateVaultRoleRequest, opts ...grpc.CallOption) (*UpdateVaultRoleResponse, error) {
	out := new(UpdateVaultRoleResponse)
	err := c.cc.Invoke(ctx, Agent_UpdateVaultRole_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *agentClient) DeleteVaultRole(ctx context.Context, in *DeleteVaultRoleRequest, opts ...grpc.CallOption) (*DeleteVaultRoleResponse, error) {
	out := new(DeleteVaultRoleResponse)
	err := c.cc.Invoke(ctx, Agent_DeleteVaultRole_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *agentClient) SyncApp(ctx context.Context, in *SyncAppRequest, opts ...grpc.CallOption) (*SyncAppResponse, error) {
	out := new(SyncAppResponse)
	err := c.cc.Invoke(ctx, Agent_SyncApp_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *agentClient) GetClusterApps(ctx context.Context, in *GetClusterAppsRequest, opts ...grpc.CallOption) (*GetClusterAppsResponse, error) {
	out := new(GetClusterAppsResponse)
	err := c.cc.Invoke(ctx, Agent_GetClusterApps_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *agentClient) GetClusterAppLaunches(ctx context.Context, in *GetClusterAppLaunchesRequest, opts ...grpc.CallOption) (*GetClusterAppLaunchesResponse, error) {
	out := new(GetClusterAppLaunchesResponse)
	err := c.cc.Invoke(ctx, Agent_GetClusterAppLaunches_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *agentClient) ConfigureAppSSO(ctx context.Context, in *ConfigureAppSSORequest, opts ...grpc.CallOption) (*ConfigureAppSSOResponse, error) {
	out := new(ConfigureAppSSOResponse)
	err := c.cc.Invoke(ctx, Agent_ConfigureAppSSO_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *agentClient) GetClusterAppConfig(ctx context.Context, in *GetClusterAppConfigRequest, opts ...grpc.CallOption) (*GetClusterAppConfigResponse, error) {
	out := new(GetClusterAppConfigResponse)
	err := c.cc.Invoke(ctx, Agent_GetClusterAppConfig_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *agentClient) GetClusterAppValues(ctx context.Context, in *GetClusterAppValuesRequest, opts ...grpc.CallOption) (*GetClusterAppValuesResponse, error) {
	out := new(GetClusterAppValuesResponse)
	err := c.cc.Invoke(ctx, Agent_GetClusterAppValues_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *agentClient) GetClusterGlobalValues(ctx context.Context, in *GetClusterGlobalValuesRequest, opts ...grpc.CallOption) (*GetClusterGlobalValuesResponse, error) {
	out := new(GetClusterGlobalValuesResponse)
	err := c.cc.Invoke(ctx, Agent_GetClusterGlobalValues_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *agentClient) DeployDefaultApps(ctx context.Context, in *DeployDefaultAppsRequest, opts ...grpc.CallOption) (*DeployDefaultAppsResponse, error) {
	out := new(DeployDefaultAppsResponse)
	err := c.cc.Invoke(ctx, Agent_DeployDefaultApps_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *agentClient) GetDefaultAppsStatus(ctx context.Context, in *GetDefaultAppsStatusRequest, opts ...grpc.CallOption) (*GetDefaultAppsStatusResponse, error) {
	out := new(GetDefaultAppsStatusResponse)
	err := c.cc.Invoke(ctx, Agent_GetDefaultAppsStatus_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AgentServer is the server API for Agent service.
// All implementations must embed UnimplementedAgentServer
// for forward compatibility
type AgentServer interface {
	Ping(context.Context, *PingRequest) (*PingResponse, error)
	StoreCredential(context.Context, *StoreCredentialRequest) (*StoreCredentialResponse, error)
	ConfigureVaultSecret(context.Context, *ConfigureVaultSecretRequest) (*ConfigureVaultSecretResponse, error)
	CreateVaultRole(context.Context, *CreateVaultRoleRequest) (*CreateVaultRoleResponse, error)
	UpdateVaultRole(context.Context, *UpdateVaultRoleRequest) (*UpdateVaultRoleResponse, error)
	DeleteVaultRole(context.Context, *DeleteVaultRoleRequest) (*DeleteVaultRoleResponse, error)
	SyncApp(context.Context, *SyncAppRequest) (*SyncAppResponse, error)
	GetClusterApps(context.Context, *GetClusterAppsRequest) (*GetClusterAppsResponse, error)
	GetClusterAppLaunches(context.Context, *GetClusterAppLaunchesRequest) (*GetClusterAppLaunchesResponse, error)
	ConfigureAppSSO(context.Context, *ConfigureAppSSORequest) (*ConfigureAppSSOResponse, error)
	GetClusterAppConfig(context.Context, *GetClusterAppConfigRequest) (*GetClusterAppConfigResponse, error)
	GetClusterAppValues(context.Context, *GetClusterAppValuesRequest) (*GetClusterAppValuesResponse, error)
	GetClusterGlobalValues(context.Context, *GetClusterGlobalValuesRequest) (*GetClusterGlobalValuesResponse, error)
	DeployDefaultApps(context.Context, *DeployDefaultAppsRequest) (*DeployDefaultAppsResponse, error)
	GetDefaultAppsStatus(context.Context, *GetDefaultAppsStatusRequest) (*GetDefaultAppsStatusResponse, error)
	mustEmbedUnimplementedAgentServer()
}

// UnimplementedAgentServer must be embedded to have forward compatible implementations.
type UnimplementedAgentServer struct {
}

func (UnimplementedAgentServer) Ping(context.Context, *PingRequest) (*PingResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Ping not implemented")
}
func (UnimplementedAgentServer) StoreCredential(context.Context, *StoreCredentialRequest) (*StoreCredentialResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StoreCredential not implemented")
}
func (UnimplementedAgentServer) ConfigureVaultSecret(context.Context, *ConfigureVaultSecretRequest) (*ConfigureVaultSecretResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ConfigureVaultSecret not implemented")
}
func (UnimplementedAgentServer) CreateVaultRole(context.Context, *CreateVaultRoleRequest) (*CreateVaultRoleResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateVaultRole not implemented")
}
func (UnimplementedAgentServer) UpdateVaultRole(context.Context, *UpdateVaultRoleRequest) (*UpdateVaultRoleResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateVaultRole not implemented")
}
func (UnimplementedAgentServer) DeleteVaultRole(context.Context, *DeleteVaultRoleRequest) (*DeleteVaultRoleResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteVaultRole not implemented")
}
func (UnimplementedAgentServer) SyncApp(context.Context, *SyncAppRequest) (*SyncAppResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SyncApp not implemented")
}
func (UnimplementedAgentServer) GetClusterApps(context.Context, *GetClusterAppsRequest) (*GetClusterAppsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetClusterApps not implemented")
}
func (UnimplementedAgentServer) GetClusterAppLaunches(context.Context, *GetClusterAppLaunchesRequest) (*GetClusterAppLaunchesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetClusterAppLaunches not implemented")
}
func (UnimplementedAgentServer) ConfigureAppSSO(context.Context, *ConfigureAppSSORequest) (*ConfigureAppSSOResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ConfigureAppSSO not implemented")
}
func (UnimplementedAgentServer) GetClusterAppConfig(context.Context, *GetClusterAppConfigRequest) (*GetClusterAppConfigResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetClusterAppConfig not implemented")
}
func (UnimplementedAgentServer) GetClusterAppValues(context.Context, *GetClusterAppValuesRequest) (*GetClusterAppValuesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetClusterAppValues not implemented")
}
func (UnimplementedAgentServer) GetClusterGlobalValues(context.Context, *GetClusterGlobalValuesRequest) (*GetClusterGlobalValuesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetClusterGlobalValues not implemented")
}
func (UnimplementedAgentServer) DeployDefaultApps(context.Context, *DeployDefaultAppsRequest) (*DeployDefaultAppsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeployDefaultApps not implemented")
}
func (UnimplementedAgentServer) GetDefaultAppsStatus(context.Context, *GetDefaultAppsStatusRequest) (*GetDefaultAppsStatusResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetDefaultAppsStatus not implemented")
}
func (UnimplementedAgentServer) mustEmbedUnimplementedAgentServer() {}

// UnsafeAgentServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AgentServer will
// result in compilation errors.
type UnsafeAgentServer interface {
	mustEmbedUnimplementedAgentServer()
}

func RegisterAgentServer(s grpc.ServiceRegistrar, srv AgentServer) {
	s.RegisterService(&Agent_ServiceDesc, srv)
}

func _Agent_Ping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AgentServer).Ping(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Agent_Ping_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AgentServer).Ping(ctx, req.(*PingRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Agent_StoreCredential_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StoreCredentialRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AgentServer).StoreCredential(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Agent_StoreCredential_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AgentServer).StoreCredential(ctx, req.(*StoreCredentialRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Agent_ConfigureVaultSecret_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ConfigureVaultSecretRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AgentServer).ConfigureVaultSecret(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Agent_ConfigureVaultSecret_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AgentServer).ConfigureVaultSecret(ctx, req.(*ConfigureVaultSecretRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Agent_CreateVaultRole_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateVaultRoleRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AgentServer).CreateVaultRole(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Agent_CreateVaultRole_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AgentServer).CreateVaultRole(ctx, req.(*CreateVaultRoleRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Agent_UpdateVaultRole_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateVaultRoleRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AgentServer).UpdateVaultRole(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Agent_UpdateVaultRole_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AgentServer).UpdateVaultRole(ctx, req.(*UpdateVaultRoleRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Agent_DeleteVaultRole_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteVaultRoleRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AgentServer).DeleteVaultRole(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Agent_DeleteVaultRole_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AgentServer).DeleteVaultRole(ctx, req.(*DeleteVaultRoleRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Agent_SyncApp_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SyncAppRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AgentServer).SyncApp(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Agent_SyncApp_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AgentServer).SyncApp(ctx, req.(*SyncAppRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Agent_GetClusterApps_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetClusterAppsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AgentServer).GetClusterApps(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Agent_GetClusterApps_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AgentServer).GetClusterApps(ctx, req.(*GetClusterAppsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Agent_GetClusterAppLaunches_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetClusterAppLaunchesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AgentServer).GetClusterAppLaunches(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Agent_GetClusterAppLaunches_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AgentServer).GetClusterAppLaunches(ctx, req.(*GetClusterAppLaunchesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Agent_ConfigureAppSSO_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ConfigureAppSSORequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AgentServer).ConfigureAppSSO(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Agent_ConfigureAppSSO_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AgentServer).ConfigureAppSSO(ctx, req.(*ConfigureAppSSORequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Agent_GetClusterAppConfig_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetClusterAppConfigRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AgentServer).GetClusterAppConfig(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Agent_GetClusterAppConfig_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AgentServer).GetClusterAppConfig(ctx, req.(*GetClusterAppConfigRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Agent_GetClusterAppValues_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetClusterAppValuesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AgentServer).GetClusterAppValues(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Agent_GetClusterAppValues_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AgentServer).GetClusterAppValues(ctx, req.(*GetClusterAppValuesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Agent_GetClusterGlobalValues_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetClusterGlobalValuesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AgentServer).GetClusterGlobalValues(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Agent_GetClusterGlobalValues_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AgentServer).GetClusterGlobalValues(ctx, req.(*GetClusterGlobalValuesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Agent_DeployDefaultApps_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeployDefaultAppsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AgentServer).DeployDefaultApps(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Agent_DeployDefaultApps_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AgentServer).DeployDefaultApps(ctx, req.(*DeployDefaultAppsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Agent_GetDefaultAppsStatus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetDefaultAppsStatusRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AgentServer).GetDefaultAppsStatus(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Agent_GetDefaultAppsStatus_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AgentServer).GetDefaultAppsStatus(ctx, req.(*GetDefaultAppsStatusRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Agent_ServiceDesc is the grpc.ServiceDesc for Agent service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Agent_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "agentpb.Agent",
	HandlerType: (*AgentServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Ping",
			Handler:    _Agent_Ping_Handler,
		},
		{
			MethodName: "StoreCredential",
			Handler:    _Agent_StoreCredential_Handler,
		},
		{
			MethodName: "ConfigureVaultSecret",
			Handler:    _Agent_ConfigureVaultSecret_Handler,
		},
		{
			MethodName: "CreateVaultRole",
			Handler:    _Agent_CreateVaultRole_Handler,
		},
		{
			MethodName: "UpdateVaultRole",
			Handler:    _Agent_UpdateVaultRole_Handler,
		},
		{
			MethodName: "DeleteVaultRole",
			Handler:    _Agent_DeleteVaultRole_Handler,
		},
		{
			MethodName: "SyncApp",
			Handler:    _Agent_SyncApp_Handler,
		},
		{
			MethodName: "GetClusterApps",
			Handler:    _Agent_GetClusterApps_Handler,
		},
		{
			MethodName: "GetClusterAppLaunches",
			Handler:    _Agent_GetClusterAppLaunches_Handler,
		},
		{
			MethodName: "ConfigureAppSSO",
			Handler:    _Agent_ConfigureAppSSO_Handler,
		},
		{
			MethodName: "GetClusterAppConfig",
			Handler:    _Agent_GetClusterAppConfig_Handler,
		},
		{
			MethodName: "GetClusterAppValues",
			Handler:    _Agent_GetClusterAppValues_Handler,
		},
		{
			MethodName: "GetClusterGlobalValues",
			Handler:    _Agent_GetClusterGlobalValues_Handler,
		},
		{
			MethodName: "DeployDefaultApps",
			Handler:    _Agent_DeployDefaultApps_Handler,
		},
		{
			MethodName: "GetDefaultAppsStatus",
			Handler:    _Agent_GetDefaultAppsStatus_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "agent.proto",
}