// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.12
// source: proto/oauthproto/oauth.proto

package oauthproto

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type OauthClientRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ClientName              string   `protobuf:"bytes,1,opt,name=client_name,json=clientName,proto3" json:"client_name,omitempty"`
	RedirectUris            []string `protobuf:"bytes,2,rep,name=redirect_uris,json=redirectUris,proto3" json:"redirect_uris,omitempty"`
	GrantTypes              []string `protobuf:"bytes,3,rep,name=grant_types,json=grantTypes,proto3" json:"grant_types,omitempty"`
	ResponseTypes           []string `protobuf:"bytes,4,rep,name=response_types,json=responseTypes,proto3" json:"response_types,omitempty"`
	TokenEndpointAuthMethod string   `protobuf:"bytes,5,opt,name=token_endpoint_auth_method,json=tokenEndpointAuthMethod,proto3" json:"token_endpoint_auth_method,omitempty"`
	Scope                   string   `protobuf:"bytes,6,opt,name=scope,proto3" json:"scope,omitempty"`
}

func (x *OauthClientRequest) Reset() {
	*x = OauthClientRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_oauthproto_oauth_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *OauthClientRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OauthClientRequest) ProtoMessage() {}

func (x *OauthClientRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_oauthproto_oauth_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OauthClientRequest.ProtoReflect.Descriptor instead.
func (*OauthClientRequest) Descriptor() ([]byte, []int) {
	return file_proto_oauthproto_oauth_proto_rawDescGZIP(), []int{0}
}

func (x *OauthClientRequest) GetClientName() string {
	if x != nil {
		return x.ClientName
	}
	return ""
}

func (x *OauthClientRequest) GetRedirectUris() []string {
	if x != nil {
		return x.RedirectUris
	}
	return nil
}

func (x *OauthClientRequest) GetGrantTypes() []string {
	if x != nil {
		return x.GrantTypes
	}
	return nil
}

func (x *OauthClientRequest) GetResponseTypes() []string {
	if x != nil {
		return x.ResponseTypes
	}
	return nil
}

func (x *OauthClientRequest) GetTokenEndpointAuthMethod() string {
	if x != nil {
		return x.TokenEndpointAuthMethod
	}
	return ""
}

func (x *OauthClientRequest) GetScope() string {
	if x != nil {
		return x.Scope
	}
	return ""
}

type OauthClientResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ClientId     string `protobuf:"bytes,1,opt,name=client_id,json=clientId,proto3" json:"client_id,omitempty"`
	ClientSecret string `protobuf:"bytes,2,opt,name=client_secret,json=clientSecret,proto3" json:"client_secret,omitempty"`
}

func (x *OauthClientResponse) Reset() {
	*x = OauthClientResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_oauthproto_oauth_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *OauthClientResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OauthClientResponse) ProtoMessage() {}

func (x *OauthClientResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_oauthproto_oauth_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OauthClientResponse.ProtoReflect.Descriptor instead.
func (*OauthClientResponse) Descriptor() ([]byte, []int) {
	return file_proto_oauthproto_oauth_proto_rawDescGZIP(), []int{1}
}

func (x *OauthClientResponse) GetClientId() string {
	if x != nil {
		return x.ClientId
	}
	return ""
}

func (x *OauthClientResponse) GetClientSecret() string {
	if x != nil {
		return x.ClientSecret
	}
	return ""
}

type OauthTokenRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ClientId     string `protobuf:"bytes,1,opt,name=client_id,json=clientId,proto3" json:"client_id,omitempty"`
	ClientSecret string `protobuf:"bytes,2,opt,name=client_secret,json=clientSecret,proto3" json:"client_secret,omitempty"`
}

func (x *OauthTokenRequest) Reset() {
	*x = OauthTokenRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_oauthproto_oauth_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *OauthTokenRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OauthTokenRequest) ProtoMessage() {}

func (x *OauthTokenRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_oauthproto_oauth_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OauthTokenRequest.ProtoReflect.Descriptor instead.
func (*OauthTokenRequest) Descriptor() ([]byte, []int) {
	return file_proto_oauthproto_oauth_proto_rawDescGZIP(), []int{2}
}

func (x *OauthTokenRequest) GetClientId() string {
	if x != nil {
		return x.ClientId
	}
	return ""
}

func (x *OauthTokenRequest) GetClientSecret() string {
	if x != nil {
		return x.ClientSecret
	}
	return ""
}

type OauthTokenResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	OauthToken   string `protobuf:"bytes,1,opt,name=oauth_token,json=oauthToken,proto3" json:"oauth_token,omitempty"`
	RefreshToken string `protobuf:"bytes,2,opt,name=refresh_token,json=refreshToken,proto3" json:"refresh_token,omitempty"`
}

func (x *OauthTokenResponse) Reset() {
	*x = OauthTokenResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_oauthproto_oauth_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *OauthTokenResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OauthTokenResponse) ProtoMessage() {}

func (x *OauthTokenResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_oauthproto_oauth_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OauthTokenResponse.ProtoReflect.Descriptor instead.
func (*OauthTokenResponse) Descriptor() ([]byte, []int) {
	return file_proto_oauthproto_oauth_proto_rawDescGZIP(), []int{3}
}

func (x *OauthTokenResponse) GetOauthToken() string {
	if x != nil {
		return x.OauthToken
	}
	return ""
}

func (x *OauthTokenResponse) GetRefreshToken() string {
	if x != nil {
		return x.RefreshToken
	}
	return ""
}

type ValidateOauthTokenRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	OauthToken string `protobuf:"bytes,1,opt,name=oauth_token,json=oauthToken,proto3" json:"oauth_token,omitempty"`
}

func (x *ValidateOauthTokenRequest) Reset() {
	*x = ValidateOauthTokenRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_oauthproto_oauth_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ValidateOauthTokenRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ValidateOauthTokenRequest) ProtoMessage() {}

func (x *ValidateOauthTokenRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_oauthproto_oauth_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ValidateOauthTokenRequest.ProtoReflect.Descriptor instead.
func (*ValidateOauthTokenRequest) Descriptor() ([]byte, []int) {
	return file_proto_oauthproto_oauth_proto_rawDescGZIP(), []int{4}
}

func (x *ValidateOauthTokenRequest) GetOauthToken() string {
	if x != nil {
		return x.OauthToken
	}
	return ""
}

type ValidateOauthTokenResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Valid string `protobuf:"bytes,1,opt,name=valid,proto3" json:"valid,omitempty"`
}

func (x *ValidateOauthTokenResponse) Reset() {
	*x = ValidateOauthTokenResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_oauthproto_oauth_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ValidateOauthTokenResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ValidateOauthTokenResponse) ProtoMessage() {}

func (x *ValidateOauthTokenResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_oauthproto_oauth_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ValidateOauthTokenResponse.ProtoReflect.Descriptor instead.
func (*ValidateOauthTokenResponse) Descriptor() ([]byte, []int) {
	return file_proto_oauthproto_oauth_proto_rawDescGZIP(), []int{5}
}

func (x *ValidateOauthTokenResponse) GetValid() string {
	if x != nil {
		return x.Valid
	}
	return ""
}

type CreateClientCredentialsClientRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ClientName string `protobuf:"bytes,1,opt,name=client_name,json=clientName,proto3" json:"client_name,omitempty"`
}

func (x *CreateClientCredentialsClientRequest) Reset() {
	*x = CreateClientCredentialsClientRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_oauthproto_oauth_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateClientCredentialsClientRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateClientCredentialsClientRequest) ProtoMessage() {}

func (x *CreateClientCredentialsClientRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_oauthproto_oauth_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateClientCredentialsClientRequest.ProtoReflect.Descriptor instead.
func (*CreateClientCredentialsClientRequest) Descriptor() ([]byte, []int) {
	return file_proto_oauthproto_oauth_proto_rawDescGZIP(), []int{6}
}

func (x *CreateClientCredentialsClientRequest) GetClientName() string {
	if x != nil {
		return x.ClientName
	}
	return ""
}

type CreateClientCredentialsClientResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ClientId     string `protobuf:"bytes,1,opt,name=client_id,json=clientId,proto3" json:"client_id,omitempty"`
	ClientSecret string `protobuf:"bytes,2,opt,name=client_secret,json=clientSecret,proto3" json:"client_secret,omitempty"`
}

func (x *CreateClientCredentialsClientResponse) Reset() {
	*x = CreateClientCredentialsClientResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_oauthproto_oauth_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateClientCredentialsClientResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateClientCredentialsClientResponse) ProtoMessage() {}

func (x *CreateClientCredentialsClientResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_oauthproto_oauth_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateClientCredentialsClientResponse.ProtoReflect.Descriptor instead.
func (*CreateClientCredentialsClientResponse) Descriptor() ([]byte, []int) {
	return file_proto_oauthproto_oauth_proto_rawDescGZIP(), []int{7}
}

func (x *CreateClientCredentialsClientResponse) GetClientId() string {
	if x != nil {
		return x.ClientId
	}
	return ""
}

func (x *CreateClientCredentialsClientResponse) GetClientSecret() string {
	if x != nil {
		return x.ClientSecret
	}
	return ""
}

var File_proto_oauthproto_oauth_proto protoreflect.FileDescriptor

var file_proto_oauthproto_oauth_proto_rawDesc = []byte{
	0x0a, 0x1c, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x6f, 0x61, 0x75, 0x74, 0x68, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2f, 0x6f, 0x61, 0x75, 0x74, 0x68, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0c,
	0x4f, 0x61, 0x75, 0x74, 0x68, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x22, 0xf5, 0x01, 0x0a,
	0x12, 0x4f, 0x61, 0x75, 0x74, 0x68, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x1f, 0x0a, 0x0b, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x5f, 0x6e, 0x61,
	0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74,
	0x4e, 0x61, 0x6d, 0x65, 0x12, 0x23, 0x0a, 0x0d, 0x72, 0x65, 0x64, 0x69, 0x72, 0x65, 0x63, 0x74,
	0x5f, 0x75, 0x72, 0x69, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x09, 0x52, 0x0c, 0x72, 0x65, 0x64,
	0x69, 0x72, 0x65, 0x63, 0x74, 0x55, 0x72, 0x69, 0x73, 0x12, 0x1f, 0x0a, 0x0b, 0x67, 0x72, 0x61,
	0x6e, 0x74, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x09, 0x52, 0x0a,
	0x67, 0x72, 0x61, 0x6e, 0x74, 0x54, 0x79, 0x70, 0x65, 0x73, 0x12, 0x25, 0x0a, 0x0e, 0x72, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x73, 0x18, 0x04, 0x20, 0x03,
	0x28, 0x09, 0x52, 0x0d, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x54, 0x79, 0x70, 0x65,
	0x73, 0x12, 0x3b, 0x0a, 0x1a, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x5f, 0x65, 0x6e, 0x64, 0x70, 0x6f,
	0x69, 0x6e, 0x74, 0x5f, 0x61, 0x75, 0x74, 0x68, 0x5f, 0x6d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x18,
	0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x17, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x45, 0x6e, 0x64, 0x70,
	0x6f, 0x69, 0x6e, 0x74, 0x41, 0x75, 0x74, 0x68, 0x4d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x12, 0x14,
	0x0a, 0x05, 0x73, 0x63, 0x6f, 0x70, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x73,
	0x63, 0x6f, 0x70, 0x65, 0x22, 0x57, 0x0a, 0x13, 0x4f, 0x61, 0x75, 0x74, 0x68, 0x43, 0x6c, 0x69,
	0x65, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1b, 0x0a, 0x09, 0x63,
	0x6c, 0x69, 0x65, 0x6e, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08,
	0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x49, 0x64, 0x12, 0x23, 0x0a, 0x0d, 0x63, 0x6c, 0x69, 0x65,
	0x6e, 0x74, 0x5f, 0x73, 0x65, 0x63, 0x72, 0x65, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0c, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x53, 0x65, 0x63, 0x72, 0x65, 0x74, 0x22, 0x55, 0x0a,
	0x11, 0x4f, 0x61, 0x75, 0x74, 0x68, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x1b, 0x0a, 0x09, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x5f, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x49, 0x64, 0x12,
	0x23, 0x0a, 0x0d, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x5f, 0x73, 0x65, 0x63, 0x72, 0x65, 0x74,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x53, 0x65,
	0x63, 0x72, 0x65, 0x74, 0x22, 0x5a, 0x0a, 0x12, 0x4f, 0x61, 0x75, 0x74, 0x68, 0x54, 0x6f, 0x6b,
	0x65, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1f, 0x0a, 0x0b, 0x6f, 0x61,
	0x75, 0x74, 0x68, 0x5f, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0a, 0x6f, 0x61, 0x75, 0x74, 0x68, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x23, 0x0a, 0x0d, 0x72,
	0x65, 0x66, 0x72, 0x65, 0x73, 0x68, 0x5f, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0c, 0x72, 0x65, 0x66, 0x72, 0x65, 0x73, 0x68, 0x54, 0x6f, 0x6b, 0x65, 0x6e,
	0x22, 0x3c, 0x0a, 0x19, 0x56, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x4f, 0x61, 0x75, 0x74,
	0x68, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1f, 0x0a,
	0x0b, 0x6f, 0x61, 0x75, 0x74, 0x68, 0x5f, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0a, 0x6f, 0x61, 0x75, 0x74, 0x68, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x22, 0x32,
	0x0a, 0x1a, 0x56, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x4f, 0x61, 0x75, 0x74, 0x68, 0x54,
	0x6f, 0x6b, 0x65, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x14, 0x0a, 0x05,
	0x76, 0x61, 0x6c, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c,
	0x69, 0x64, 0x22, 0x47, 0x0a, 0x24, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x43, 0x6c, 0x69, 0x65,
	0x6e, 0x74, 0x43, 0x72, 0x65, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x61, 0x6c, 0x73, 0x43, 0x6c, 0x69,
	0x65, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1f, 0x0a, 0x0b, 0x63, 0x6c,
	0x69, 0x65, 0x6e, 0x74, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0a, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x4e, 0x61, 0x6d, 0x65, 0x22, 0x69, 0x0a, 0x25, 0x43,
	0x72, 0x65, 0x61, 0x74, 0x65, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x43, 0x72, 0x65, 0x64, 0x65,
	0x6e, 0x74, 0x69, 0x61, 0x6c, 0x73, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1b, 0x0a, 0x09, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x5f, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x49,
	0x64, 0x12, 0x23, 0x0a, 0x0d, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x5f, 0x73, 0x65, 0x63, 0x72,
	0x65, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74,
	0x53, 0x65, 0x63, 0x72, 0x65, 0x74, 0x32, 0xb0, 0x03, 0x0a, 0x0c, 0x4f, 0x61, 0x75, 0x74, 0x68,
	0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x58, 0x0a, 0x11, 0x43, 0x72, 0x65, 0x61, 0x74,
	0x65, 0x4f, 0x61, 0x75, 0x74, 0x68, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x12, 0x20, 0x2e, 0x4f,
	0x61, 0x75, 0x74, 0x68, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x4f, 0x61, 0x75, 0x74,
	0x68, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x21,
	0x2e, 0x4f, 0x61, 0x75, 0x74, 0x68, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x4f, 0x61,
	0x75, 0x74, 0x68, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x52, 0x0a, 0x0d, 0x47, 0x65, 0x74, 0x4f, 0x61, 0x75, 0x74, 0x68, 0x54, 0x6f, 0x6b,
	0x65, 0x6e, 0x12, 0x1f, 0x2e, 0x4f, 0x61, 0x75, 0x74, 0x68, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x2e, 0x4f, 0x61, 0x75, 0x74, 0x68, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x20, 0x2e, 0x4f, 0x61, 0x75, 0x74, 0x68, 0x53, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x2e, 0x4f, 0x61, 0x75, 0x74, 0x68, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x67, 0x0a, 0x12, 0x56, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74,
	0x65, 0x4f, 0x61, 0x75, 0x74, 0x68, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x27, 0x2e, 0x4f, 0x61,
	0x75, 0x74, 0x68, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x56, 0x61, 0x6c, 0x69, 0x64,
	0x61, 0x74, 0x65, 0x4f, 0x61, 0x75, 0x74, 0x68, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x28, 0x2e, 0x4f, 0x61, 0x75, 0x74, 0x68, 0x53, 0x65, 0x72, 0x76,
	0x69, 0x63, 0x65, 0x2e, 0x56, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x4f, 0x61, 0x75, 0x74,
	0x68, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x88,
	0x01, 0x0a, 0x1d, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x43,
	0x72, 0x65, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x61, 0x6c, 0x73, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74,
	0x12, 0x32, 0x2e, 0x4f, 0x61, 0x75, 0x74, 0x68, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e,
	0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x43, 0x72, 0x65, 0x64,
	0x65, 0x6e, 0x74, 0x69, 0x61, 0x6c, 0x73, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x33, 0x2e, 0x4f, 0x61, 0x75, 0x74, 0x68, 0x53, 0x65, 0x72, 0x76,
	0x69, 0x63, 0x65, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74,
	0x43, 0x72, 0x65, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x61, 0x6c, 0x73, 0x43, 0x6c, 0x69, 0x65, 0x6e,
	0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x2e, 0x5a, 0x2c, 0x67, 0x69, 0x74,
	0x6c, 0x61, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x74, 0x61, 0x72, 0x69, 0x61, 0x6e, 0x64, 0x65,
	0x76, 0x5f, 0x69, 0x6e, 0x74, 0x65, 0x6c, 0x6f, 0x70, 0x73, 0x2f, 0x69, 0x61, 0x6d, 0x2f, 0x6f,
	0x61, 0x75, 0x74, 0x68, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_proto_oauthproto_oauth_proto_rawDescOnce sync.Once
	file_proto_oauthproto_oauth_proto_rawDescData = file_proto_oauthproto_oauth_proto_rawDesc
)

func file_proto_oauthproto_oauth_proto_rawDescGZIP() []byte {
	file_proto_oauthproto_oauth_proto_rawDescOnce.Do(func() {
		file_proto_oauthproto_oauth_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_oauthproto_oauth_proto_rawDescData)
	})
	return file_proto_oauthproto_oauth_proto_rawDescData
}

var file_proto_oauthproto_oauth_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_proto_oauthproto_oauth_proto_goTypes = []interface{}{
	(*OauthClientRequest)(nil),                    // 0: OauthService.OauthClientRequest
	(*OauthClientResponse)(nil),                   // 1: OauthService.OauthClientResponse
	(*OauthTokenRequest)(nil),                     // 2: OauthService.OauthTokenRequest
	(*OauthTokenResponse)(nil),                    // 3: OauthService.OauthTokenResponse
	(*ValidateOauthTokenRequest)(nil),             // 4: OauthService.ValidateOauthTokenRequest
	(*ValidateOauthTokenResponse)(nil),            // 5: OauthService.ValidateOauthTokenResponse
	(*CreateClientCredentialsClientRequest)(nil),  // 6: OauthService.CreateClientCredentialsClientRequest
	(*CreateClientCredentialsClientResponse)(nil), // 7: OauthService.CreateClientCredentialsClientResponse
}
var file_proto_oauthproto_oauth_proto_depIdxs = []int32{
	0, // 0: OauthService.OauthService.CreateOauthClient:input_type -> OauthService.OauthClientRequest
	2, // 1: OauthService.OauthService.GetOauthToken:input_type -> OauthService.OauthTokenRequest
	4, // 2: OauthService.OauthService.ValidateOauthToken:input_type -> OauthService.ValidateOauthTokenRequest
	6, // 3: OauthService.OauthService.CreateClientCredentialsClient:input_type -> OauthService.CreateClientCredentialsClientRequest
	1, // 4: OauthService.OauthService.CreateOauthClient:output_type -> OauthService.OauthClientResponse
	3, // 5: OauthService.OauthService.GetOauthToken:output_type -> OauthService.OauthTokenResponse
	5, // 6: OauthService.OauthService.ValidateOauthToken:output_type -> OauthService.ValidateOauthTokenResponse
	7, // 7: OauthService.OauthService.CreateClientCredentialsClient:output_type -> OauthService.CreateClientCredentialsClientResponse
	4, // [4:8] is the sub-list for method output_type
	0, // [0:4] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_proto_oauthproto_oauth_proto_init() }
func file_proto_oauthproto_oauth_proto_init() {
	if File_proto_oauthproto_oauth_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_oauthproto_oauth_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*OauthClientRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_oauthproto_oauth_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*OauthClientResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_oauthproto_oauth_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*OauthTokenRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_oauthproto_oauth_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*OauthTokenResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_oauthproto_oauth_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ValidateOauthTokenRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_oauthproto_oauth_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ValidateOauthTokenResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_oauthproto_oauth_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateClientCredentialsClientRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_oauthproto_oauth_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateClientCredentialsClientResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_proto_oauthproto_oauth_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_oauthproto_oauth_proto_goTypes,
		DependencyIndexes: file_proto_oauthproto_oauth_proto_depIdxs,
		MessageInfos:      file_proto_oauthproto_oauth_proto_msgTypes,
	}.Build()
	File_proto_oauthproto_oauth_proto = out.File
	file_proto_oauthproto_oauth_proto_rawDesc = nil
	file_proto_oauthproto_oauth_proto_goTypes = nil
	file_proto_oauthproto_oauth_proto_depIdxs = nil
}
