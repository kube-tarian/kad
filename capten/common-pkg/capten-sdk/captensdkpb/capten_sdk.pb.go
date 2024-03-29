// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.30.0
// 	protoc        v3.12.4
// source: capten_sdk.proto

package captensdkpb

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

type StatusCode int32

const (
	StatusCode_OK               StatusCode = 0
	StatusCode_INTERNAL_ERROR   StatusCode = 1
	StatusCode_INVALID_ARGUMENT StatusCode = 2
	StatusCode_NOT_FOUND        StatusCode = 3
)

// Enum value maps for StatusCode.
var (
	StatusCode_name = map[int32]string{
		0: "OK",
		1: "INTERNAL_ERROR",
		2: "INVALID_ARGUMENT",
		3: "NOT_FOUND",
	}
	StatusCode_value = map[string]int32{
		"OK":               0,
		"INTERNAL_ERROR":   1,
		"INVALID_ARGUMENT": 2,
		"NOT_FOUND":        3,
	}
)

func (x StatusCode) Enum() *StatusCode {
	p := new(StatusCode)
	*p = x
	return p
}

func (x StatusCode) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (StatusCode) Descriptor() protoreflect.EnumDescriptor {
	return file_capten_sdk_proto_enumTypes[0].Descriptor()
}

func (StatusCode) Type() protoreflect.EnumType {
	return &file_capten_sdk_proto_enumTypes[0]
}

func (x StatusCode) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use StatusCode.Descriptor instead.
func (StatusCode) EnumDescriptor() ([]byte, []int) {
	return file_capten_sdk_proto_rawDescGZIP(), []int{0}
}

type GitProject struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id             string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	ProjectUrl     string   `protobuf:"bytes,2,opt,name=projectUrl,proto3" json:"projectUrl,omitempty"`
	AccessToken    string   `protobuf:"bytes,3,opt,name=accessToken,proto3" json:"accessToken,omitempty"`
	Labels         []string `protobuf:"bytes,4,rep,name=labels,proto3" json:"labels,omitempty"`
	LastUpdateTime string   `protobuf:"bytes,5,opt,name=lastUpdateTime,proto3" json:"lastUpdateTime,omitempty"`
}

func (x *GitProject) Reset() {
	*x = GitProject{}
	if protoimpl.UnsafeEnabled {
		mi := &file_capten_sdk_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GitProject) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GitProject) ProtoMessage() {}

func (x *GitProject) ProtoReflect() protoreflect.Message {
	mi := &file_capten_sdk_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GitProject.ProtoReflect.Descriptor instead.
func (*GitProject) Descriptor() ([]byte, []int) {
	return file_capten_sdk_proto_rawDescGZIP(), []int{0}
}

func (x *GitProject) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *GitProject) GetProjectUrl() string {
	if x != nil {
		return x.ProjectUrl
	}
	return ""
}

func (x *GitProject) GetAccessToken() string {
	if x != nil {
		return x.AccessToken
	}
	return ""
}

func (x *GitProject) GetLabels() []string {
	if x != nil {
		return x.Labels
	}
	return nil
}

func (x *GitProject) GetLastUpdateTime() string {
	if x != nil {
		return x.LastUpdateTime
	}
	return ""
}

type GetGitProjectByIdRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *GetGitProjectByIdRequest) Reset() {
	*x = GetGitProjectByIdRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_capten_sdk_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetGitProjectByIdRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetGitProjectByIdRequest) ProtoMessage() {}

func (x *GetGitProjectByIdRequest) ProtoReflect() protoreflect.Message {
	mi := &file_capten_sdk_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetGitProjectByIdRequest.ProtoReflect.Descriptor instead.
func (*GetGitProjectByIdRequest) Descriptor() ([]byte, []int) {
	return file_capten_sdk_proto_rawDescGZIP(), []int{1}
}

func (x *GetGitProjectByIdRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type GetGitProjectByIdResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Project       *GitProject `protobuf:"bytes,1,opt,name=project,proto3" json:"project,omitempty"`
	Status        StatusCode  `protobuf:"varint,2,opt,name=status,proto3,enum=captensdkpb.StatusCode" json:"status,omitempty"`
	StatusMessage string      `protobuf:"bytes,3,opt,name=statusMessage,proto3" json:"statusMessage,omitempty"`
}

func (x *GetGitProjectByIdResponse) Reset() {
	*x = GetGitProjectByIdResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_capten_sdk_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetGitProjectByIdResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetGitProjectByIdResponse) ProtoMessage() {}

func (x *GetGitProjectByIdResponse) ProtoReflect() protoreflect.Message {
	mi := &file_capten_sdk_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetGitProjectByIdResponse.ProtoReflect.Descriptor instead.
func (*GetGitProjectByIdResponse) Descriptor() ([]byte, []int) {
	return file_capten_sdk_proto_rawDescGZIP(), []int{2}
}

func (x *GetGitProjectByIdResponse) GetProject() *GitProject {
	if x != nil {
		return x.Project
	}
	return nil
}

func (x *GetGitProjectByIdResponse) GetStatus() StatusCode {
	if x != nil {
		return x.Status
	}
	return StatusCode_OK
}

func (x *GetGitProjectByIdResponse) GetStatusMessage() string {
	if x != nil {
		return x.StatusMessage
	}
	return ""
}

type ContainerRegistry struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id                 string            `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	RegistryUrl        string            `protobuf:"bytes,2,opt,name=registryUrl,proto3" json:"registryUrl,omitempty"`
	Labels             []string          `protobuf:"bytes,3,rep,name=labels,proto3" json:"labels,omitempty"`
	LastUpdateTime     string            `protobuf:"bytes,4,opt,name=lastUpdateTime,proto3" json:"lastUpdateTime,omitempty"`
	RegistryAttributes map[string]string `protobuf:"bytes,5,rep,name=registryAttributes,proto3" json:"registryAttributes,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	RegistryType       string            `protobuf:"bytes,6,opt,name=registryType,proto3" json:"registryType,omitempty"`
}

func (x *ContainerRegistry) Reset() {
	*x = ContainerRegistry{}
	if protoimpl.UnsafeEnabled {
		mi := &file_capten_sdk_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ContainerRegistry) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ContainerRegistry) ProtoMessage() {}

func (x *ContainerRegistry) ProtoReflect() protoreflect.Message {
	mi := &file_capten_sdk_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ContainerRegistry.ProtoReflect.Descriptor instead.
func (*ContainerRegistry) Descriptor() ([]byte, []int) {
	return file_capten_sdk_proto_rawDescGZIP(), []int{3}
}

func (x *ContainerRegistry) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *ContainerRegistry) GetRegistryUrl() string {
	if x != nil {
		return x.RegistryUrl
	}
	return ""
}

func (x *ContainerRegistry) GetLabels() []string {
	if x != nil {
		return x.Labels
	}
	return nil
}

func (x *ContainerRegistry) GetLastUpdateTime() string {
	if x != nil {
		return x.LastUpdateTime
	}
	return ""
}

func (x *ContainerRegistry) GetRegistryAttributes() map[string]string {
	if x != nil {
		return x.RegistryAttributes
	}
	return nil
}

func (x *ContainerRegistry) GetRegistryType() string {
	if x != nil {
		return x.RegistryType
	}
	return ""
}

type GetContainerRegistryByIdRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *GetContainerRegistryByIdRequest) Reset() {
	*x = GetContainerRegistryByIdRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_capten_sdk_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetContainerRegistryByIdRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetContainerRegistryByIdRequest) ProtoMessage() {}

func (x *GetContainerRegistryByIdRequest) ProtoReflect() protoreflect.Message {
	mi := &file_capten_sdk_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetContainerRegistryByIdRequest.ProtoReflect.Descriptor instead.
func (*GetContainerRegistryByIdRequest) Descriptor() ([]byte, []int) {
	return file_capten_sdk_proto_rawDescGZIP(), []int{4}
}

func (x *GetContainerRegistryByIdRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type GetContainerRegistryByIdResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Registry      *ContainerRegistry `protobuf:"bytes,1,opt,name=registry,proto3" json:"registry,omitempty"`
	Status        StatusCode         `protobuf:"varint,2,opt,name=status,proto3,enum=captensdkpb.StatusCode" json:"status,omitempty"`
	StatusMessage string             `protobuf:"bytes,3,opt,name=statusMessage,proto3" json:"statusMessage,omitempty"`
}

func (x *GetContainerRegistryByIdResponse) Reset() {
	*x = GetContainerRegistryByIdResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_capten_sdk_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetContainerRegistryByIdResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetContainerRegistryByIdResponse) ProtoMessage() {}

func (x *GetContainerRegistryByIdResponse) ProtoReflect() protoreflect.Message {
	mi := &file_capten_sdk_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetContainerRegistryByIdResponse.ProtoReflect.Descriptor instead.
func (*GetContainerRegistryByIdResponse) Descriptor() ([]byte, []int) {
	return file_capten_sdk_proto_rawDescGZIP(), []int{5}
}

func (x *GetContainerRegistryByIdResponse) GetRegistry() *ContainerRegistry {
	if x != nil {
		return x.Registry
	}
	return nil
}

func (x *GetContainerRegistryByIdResponse) GetStatus() StatusCode {
	if x != nil {
		return x.Status
	}
	return StatusCode_OK
}

func (x *GetContainerRegistryByIdResponse) GetStatusMessage() string {
	if x != nil {
		return x.StatusMessage
	}
	return ""
}

type DBSetupRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PluginName      string `protobuf:"bytes,1,opt,name=pluginName,proto3" json:"pluginName,omitempty"`
	DbOemName       string `protobuf:"bytes,2,opt,name=dbOemName,proto3" json:"dbOemName,omitempty"`
	DbName          string `protobuf:"bytes,3,opt,name=dbName,proto3" json:"dbName,omitempty"`
	ServiceUserName string `protobuf:"bytes,4,opt,name=serviceUserName,proto3" json:"serviceUserName,omitempty"`
}

func (x *DBSetupRequest) Reset() {
	*x = DBSetupRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_capten_sdk_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DBSetupRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DBSetupRequest) ProtoMessage() {}

func (x *DBSetupRequest) ProtoReflect() protoreflect.Message {
	mi := &file_capten_sdk_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DBSetupRequest.ProtoReflect.Descriptor instead.
func (*DBSetupRequest) Descriptor() ([]byte, []int) {
	return file_capten_sdk_proto_rawDescGZIP(), []int{6}
}

func (x *DBSetupRequest) GetPluginName() string {
	if x != nil {
		return x.PluginName
	}
	return ""
}

func (x *DBSetupRequest) GetDbOemName() string {
	if x != nil {
		return x.DbOemName
	}
	return ""
}

func (x *DBSetupRequest) GetDbName() string {
	if x != nil {
		return x.DbName
	}
	return ""
}

func (x *DBSetupRequest) GetServiceUserName() string {
	if x != nil {
		return x.ServiceUserName
	}
	return ""
}

type DBSetupResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Status        StatusCode `protobuf:"varint,1,opt,name=status,proto3,enum=captensdkpb.StatusCode" json:"status,omitempty"`
	StatusMessage string     `protobuf:"bytes,2,opt,name=statusMessage,proto3" json:"statusMessage,omitempty"`
	VaultPath     string     `protobuf:"bytes,3,opt,name=vaultPath,proto3" json:"vaultPath,omitempty"`
}

func (x *DBSetupResponse) Reset() {
	*x = DBSetupResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_capten_sdk_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DBSetupResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DBSetupResponse) ProtoMessage() {}

func (x *DBSetupResponse) ProtoReflect() protoreflect.Message {
	mi := &file_capten_sdk_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DBSetupResponse.ProtoReflect.Descriptor instead.
func (*DBSetupResponse) Descriptor() ([]byte, []int) {
	return file_capten_sdk_proto_rawDescGZIP(), []int{7}
}

func (x *DBSetupResponse) GetStatus() StatusCode {
	if x != nil {
		return x.Status
	}
	return StatusCode_OK
}

func (x *DBSetupResponse) GetStatusMessage() string {
	if x != nil {
		return x.StatusMessage
	}
	return ""
}

func (x *DBSetupResponse) GetVaultPath() string {
	if x != nil {
		return x.VaultPath
	}
	return ""
}

var File_capten_sdk_proto protoreflect.FileDescriptor

var file_capten_sdk_proto_rawDesc = []byte{
	0x0a, 0x10, 0x63, 0x61, 0x70, 0x74, 0x65, 0x6e, 0x5f, 0x73, 0x64, 0x6b, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x0b, 0x63, 0x61, 0x70, 0x74, 0x65, 0x6e, 0x73, 0x64, 0x6b, 0x70, 0x62, 0x22,
	0x9e, 0x01, 0x0a, 0x0a, 0x47, 0x69, 0x74, 0x50, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x12, 0x0e,
	0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x1e,
	0x0a, 0x0a, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x55, 0x72, 0x6c, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0a, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x55, 0x72, 0x6c, 0x12, 0x20,
	0x0a, 0x0b, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0b, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x54, 0x6f, 0x6b, 0x65, 0x6e,
	0x12, 0x16, 0x0a, 0x06, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x73, 0x18, 0x04, 0x20, 0x03, 0x28, 0x09,
	0x52, 0x06, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x73, 0x12, 0x26, 0x0a, 0x0e, 0x6c, 0x61, 0x73, 0x74,
	0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x54, 0x69, 0x6d, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0e, 0x6c, 0x61, 0x73, 0x74, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x54, 0x69, 0x6d, 0x65,
	0x22, 0x2a, 0x0a, 0x18, 0x47, 0x65, 0x74, 0x47, 0x69, 0x74, 0x50, 0x72, 0x6f, 0x6a, 0x65, 0x63,
	0x74, 0x42, 0x79, 0x49, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02,
	0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x22, 0xa5, 0x01, 0x0a,
	0x19, 0x47, 0x65, 0x74, 0x47, 0x69, 0x74, 0x50, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x42, 0x79,
	0x49, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x31, 0x0a, 0x07, 0x70, 0x72,
	0x6f, 0x6a, 0x65, 0x63, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x63, 0x61,
	0x70, 0x74, 0x65, 0x6e, 0x73, 0x64, 0x6b, 0x70, 0x62, 0x2e, 0x47, 0x69, 0x74, 0x50, 0x72, 0x6f,
	0x6a, 0x65, 0x63, 0x74, 0x52, 0x07, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x12, 0x2f, 0x0a,
	0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x17, 0x2e,
	0x63, 0x61, 0x70, 0x74, 0x65, 0x6e, 0x73, 0x64, 0x6b, 0x70, 0x62, 0x2e, 0x53, 0x74, 0x61, 0x74,
	0x75, 0x73, 0x43, 0x6f, 0x64, 0x65, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x24,
	0x0a, 0x0d, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x4d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x22, 0xd8, 0x02, 0x0a, 0x11, 0x43, 0x6f, 0x6e, 0x74, 0x61, 0x69, 0x6e,
	0x65, 0x72, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x79, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x20, 0x0a, 0x0b, 0x72, 0x65,
	0x67, 0x69, 0x73, 0x74, 0x72, 0x79, 0x55, 0x72, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0b, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x79, 0x55, 0x72, 0x6c, 0x12, 0x16, 0x0a, 0x06,
	0x6c, 0x61, 0x62, 0x65, 0x6c, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x09, 0x52, 0x06, 0x6c, 0x61,
	0x62, 0x65, 0x6c, 0x73, 0x12, 0x26, 0x0a, 0x0e, 0x6c, 0x61, 0x73, 0x74, 0x55, 0x70, 0x64, 0x61,
	0x74, 0x65, 0x54, 0x69, 0x6d, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x6c, 0x61,
	0x73, 0x74, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x66, 0x0a, 0x12,
	0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x79, 0x41, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74,
	0x65, 0x73, 0x18, 0x05, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x36, 0x2e, 0x63, 0x61, 0x70, 0x74, 0x65,
	0x6e, 0x73, 0x64, 0x6b, 0x70, 0x62, 0x2e, 0x43, 0x6f, 0x6e, 0x74, 0x61, 0x69, 0x6e, 0x65, 0x72,
	0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x79, 0x2e, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72,
	0x79, 0x41, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79,
	0x52, 0x12, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x79, 0x41, 0x74, 0x74, 0x72, 0x69, 0x62,
	0x75, 0x74, 0x65, 0x73, 0x12, 0x22, 0x0a, 0x0c, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x79,
	0x54, 0x79, 0x70, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x72, 0x65, 0x67, 0x69,
	0x73, 0x74, 0x72, 0x79, 0x54, 0x79, 0x70, 0x65, 0x1a, 0x45, 0x0a, 0x17, 0x52, 0x65, 0x67, 0x69,
	0x73, 0x74, 0x72, 0x79, 0x41, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x73, 0x45, 0x6e,
	0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x22,
	0x31, 0x0a, 0x1f, 0x47, 0x65, 0x74, 0x43, 0x6f, 0x6e, 0x74, 0x61, 0x69, 0x6e, 0x65, 0x72, 0x52,
	0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x79, 0x42, 0x79, 0x49, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02,
	0x69, 0x64, 0x22, 0xb5, 0x01, 0x0a, 0x20, 0x47, 0x65, 0x74, 0x43, 0x6f, 0x6e, 0x74, 0x61, 0x69,
	0x6e, 0x65, 0x72, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x79, 0x42, 0x79, 0x49, 0x64, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x3a, 0x0a, 0x08, 0x72, 0x65, 0x67, 0x69, 0x73,
	0x74, 0x72, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x63, 0x61, 0x70, 0x74,
	0x65, 0x6e, 0x73, 0x64, 0x6b, 0x70, 0x62, 0x2e, 0x43, 0x6f, 0x6e, 0x74, 0x61, 0x69, 0x6e, 0x65,
	0x72, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x79, 0x52, 0x08, 0x72, 0x65, 0x67, 0x69, 0x73,
	0x74, 0x72, 0x79, 0x12, 0x2f, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x0e, 0x32, 0x17, 0x2e, 0x63, 0x61, 0x70, 0x74, 0x65, 0x6e, 0x73, 0x64, 0x6b, 0x70,
	0x62, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x43, 0x6f, 0x64, 0x65, 0x52, 0x06, 0x73, 0x74,
	0x61, 0x74, 0x75, 0x73, 0x12, 0x24, 0x0a, 0x0d, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x4d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x73, 0x74, 0x61,
	0x74, 0x75, 0x73, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x22, 0x90, 0x01, 0x0a, 0x0e, 0x44,
	0x42, 0x53, 0x65, 0x74, 0x75, 0x70, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1e, 0x0a,
	0x0a, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0a, 0x70, 0x6c, 0x75, 0x67, 0x69, 0x6e, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x1c, 0x0a,
	0x09, 0x64, 0x62, 0x4f, 0x65, 0x6d, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x09, 0x64, 0x62, 0x4f, 0x65, 0x6d, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x64,
	0x62, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x64, 0x62, 0x4e,
	0x61, 0x6d, 0x65, 0x12, 0x28, 0x0a, 0x0f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x55, 0x73,
	0x65, 0x72, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0f, 0x73, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x55, 0x73, 0x65, 0x72, 0x4e, 0x61, 0x6d, 0x65, 0x22, 0x86, 0x01,
	0x0a, 0x0f, 0x44, 0x42, 0x53, 0x65, 0x74, 0x75, 0x70, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x2f, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0e, 0x32, 0x17, 0x2e, 0x63, 0x61, 0x70, 0x74, 0x65, 0x6e, 0x73, 0x64, 0x6b, 0x70, 0x62, 0x2e,
	0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x43, 0x6f, 0x64, 0x65, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74,
	0x75, 0x73, 0x12, 0x24, 0x0a, 0x0d, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x4d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x73, 0x74, 0x61, 0x74, 0x75,
	0x73, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x76, 0x61, 0x75, 0x6c,
	0x74, 0x50, 0x61, 0x74, 0x68, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x76, 0x61, 0x75,
	0x6c, 0x74, 0x50, 0x61, 0x74, 0x68, 0x2a, 0x4d, 0x0a, 0x0a, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73,
	0x43, 0x6f, 0x64, 0x65, 0x12, 0x06, 0x0a, 0x02, 0x4f, 0x4b, 0x10, 0x00, 0x12, 0x12, 0x0a, 0x0e,
	0x49, 0x4e, 0x54, 0x45, 0x52, 0x4e, 0x41, 0x4c, 0x5f, 0x45, 0x52, 0x52, 0x4f, 0x52, 0x10, 0x01,
	0x12, 0x14, 0x0a, 0x10, 0x49, 0x4e, 0x56, 0x41, 0x4c, 0x49, 0x44, 0x5f, 0x41, 0x52, 0x47, 0x55,
	0x4d, 0x45, 0x4e, 0x54, 0x10, 0x02, 0x12, 0x0d, 0x0a, 0x09, 0x4e, 0x4f, 0x54, 0x5f, 0x46, 0x4f,
	0x55, 0x4e, 0x44, 0x10, 0x03, 0x32, 0xbb, 0x02, 0x0a, 0x0a, 0x63, 0x61, 0x70, 0x74, 0x65, 0x6e,
	0x5f, 0x73, 0x64, 0x6b, 0x12, 0x64, 0x0a, 0x11, 0x47, 0x65, 0x74, 0x47, 0x69, 0x74, 0x50, 0x72,
	0x6f, 0x6a, 0x65, 0x63, 0x74, 0x42, 0x79, 0x49, 0x64, 0x12, 0x25, 0x2e, 0x63, 0x61, 0x70, 0x74,
	0x65, 0x6e, 0x73, 0x64, 0x6b, 0x70, 0x62, 0x2e, 0x47, 0x65, 0x74, 0x47, 0x69, 0x74, 0x50, 0x72,
	0x6f, 0x6a, 0x65, 0x63, 0x74, 0x42, 0x79, 0x49, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x26, 0x2e, 0x63, 0x61, 0x70, 0x74, 0x65, 0x6e, 0x73, 0x64, 0x6b, 0x70, 0x62, 0x2e, 0x47,
	0x65, 0x74, 0x47, 0x69, 0x74, 0x50, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x42, 0x79, 0x49, 0x64,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x79, 0x0a, 0x18, 0x47, 0x65,
	0x74, 0x43, 0x6f, 0x6e, 0x74, 0x61, 0x69, 0x6e, 0x65, 0x72, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74,
	0x72, 0x79, 0x42, 0x79, 0x49, 0x64, 0x12, 0x2c, 0x2e, 0x63, 0x61, 0x70, 0x74, 0x65, 0x6e, 0x73,
	0x64, 0x6b, 0x70, 0x62, 0x2e, 0x47, 0x65, 0x74, 0x43, 0x6f, 0x6e, 0x74, 0x61, 0x69, 0x6e, 0x65,
	0x72, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x79, 0x42, 0x79, 0x49, 0x64, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x2d, 0x2e, 0x63, 0x61, 0x70, 0x74, 0x65, 0x6e, 0x73, 0x64, 0x6b,
	0x70, 0x62, 0x2e, 0x47, 0x65, 0x74, 0x43, 0x6f, 0x6e, 0x74, 0x61, 0x69, 0x6e, 0x65, 0x72, 0x52,
	0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x79, 0x42, 0x79, 0x49, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x4c, 0x0a, 0x0d, 0x53, 0x65, 0x74, 0x75, 0x70, 0x44, 0x61,
	0x74, 0x61, 0x62, 0x61, 0x73, 0x65, 0x12, 0x1b, 0x2e, 0x63, 0x61, 0x70, 0x74, 0x65, 0x6e, 0x73,
	0x64, 0x6b, 0x70, 0x62, 0x2e, 0x44, 0x42, 0x53, 0x65, 0x74, 0x75, 0x70, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x1c, 0x2e, 0x63, 0x61, 0x70, 0x74, 0x65, 0x6e, 0x73, 0x64, 0x6b, 0x70,
	0x62, 0x2e, 0x44, 0x42, 0x53, 0x65, 0x74, 0x75, 0x70, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x22, 0x00, 0x42, 0x0e, 0x5a, 0x0c, 0x2f, 0x63, 0x61, 0x70, 0x74, 0x65, 0x6e, 0x73, 0x64,
	0x6b, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_capten_sdk_proto_rawDescOnce sync.Once
	file_capten_sdk_proto_rawDescData = file_capten_sdk_proto_rawDesc
)

func file_capten_sdk_proto_rawDescGZIP() []byte {
	file_capten_sdk_proto_rawDescOnce.Do(func() {
		file_capten_sdk_proto_rawDescData = protoimpl.X.CompressGZIP(file_capten_sdk_proto_rawDescData)
	})
	return file_capten_sdk_proto_rawDescData
}

var file_capten_sdk_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_capten_sdk_proto_msgTypes = make([]protoimpl.MessageInfo, 9)
var file_capten_sdk_proto_goTypes = []interface{}{
	(StatusCode)(0),                          // 0: captensdkpb.StatusCode
	(*GitProject)(nil),                       // 1: captensdkpb.GitProject
	(*GetGitProjectByIdRequest)(nil),         // 2: captensdkpb.GetGitProjectByIdRequest
	(*GetGitProjectByIdResponse)(nil),        // 3: captensdkpb.GetGitProjectByIdResponse
	(*ContainerRegistry)(nil),                // 4: captensdkpb.ContainerRegistry
	(*GetContainerRegistryByIdRequest)(nil),  // 5: captensdkpb.GetContainerRegistryByIdRequest
	(*GetContainerRegistryByIdResponse)(nil), // 6: captensdkpb.GetContainerRegistryByIdResponse
	(*DBSetupRequest)(nil),                   // 7: captensdkpb.DBSetupRequest
	(*DBSetupResponse)(nil),                  // 8: captensdkpb.DBSetupResponse
	nil,                                      // 9: captensdkpb.ContainerRegistry.RegistryAttributesEntry
}
var file_capten_sdk_proto_depIdxs = []int32{
	1, // 0: captensdkpb.GetGitProjectByIdResponse.project:type_name -> captensdkpb.GitProject
	0, // 1: captensdkpb.GetGitProjectByIdResponse.status:type_name -> captensdkpb.StatusCode
	9, // 2: captensdkpb.ContainerRegistry.registryAttributes:type_name -> captensdkpb.ContainerRegistry.RegistryAttributesEntry
	4, // 3: captensdkpb.GetContainerRegistryByIdResponse.registry:type_name -> captensdkpb.ContainerRegistry
	0, // 4: captensdkpb.GetContainerRegistryByIdResponse.status:type_name -> captensdkpb.StatusCode
	0, // 5: captensdkpb.DBSetupResponse.status:type_name -> captensdkpb.StatusCode
	2, // 6: captensdkpb.capten_sdk.GetGitProjectById:input_type -> captensdkpb.GetGitProjectByIdRequest
	5, // 7: captensdkpb.capten_sdk.GetContainerRegistryById:input_type -> captensdkpb.GetContainerRegistryByIdRequest
	7, // 8: captensdkpb.capten_sdk.SetupDatabase:input_type -> captensdkpb.DBSetupRequest
	3, // 9: captensdkpb.capten_sdk.GetGitProjectById:output_type -> captensdkpb.GetGitProjectByIdResponse
	6, // 10: captensdkpb.capten_sdk.GetContainerRegistryById:output_type -> captensdkpb.GetContainerRegistryByIdResponse
	8, // 11: captensdkpb.capten_sdk.SetupDatabase:output_type -> captensdkpb.DBSetupResponse
	9, // [9:12] is the sub-list for method output_type
	6, // [6:9] is the sub-list for method input_type
	6, // [6:6] is the sub-list for extension type_name
	6, // [6:6] is the sub-list for extension extendee
	0, // [0:6] is the sub-list for field type_name
}

func init() { file_capten_sdk_proto_init() }
func file_capten_sdk_proto_init() {
	if File_capten_sdk_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_capten_sdk_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GitProject); i {
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
		file_capten_sdk_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetGitProjectByIdRequest); i {
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
		file_capten_sdk_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetGitProjectByIdResponse); i {
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
		file_capten_sdk_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ContainerRegistry); i {
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
		file_capten_sdk_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetContainerRegistryByIdRequest); i {
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
		file_capten_sdk_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetContainerRegistryByIdResponse); i {
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
		file_capten_sdk_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DBSetupRequest); i {
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
		file_capten_sdk_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DBSetupResponse); i {
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
			RawDescriptor: file_capten_sdk_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   9,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_capten_sdk_proto_goTypes,
		DependencyIndexes: file_capten_sdk_proto_depIdxs,
		EnumInfos:         file_capten_sdk_proto_enumTypes,
		MessageInfos:      file_capten_sdk_proto_msgTypes,
	}.Build()
	File_capten_sdk_proto = out.File
	file_capten_sdk_proto_rawDesc = nil
	file_capten_sdk_proto_goTypes = nil
	file_capten_sdk_proto_depIdxs = nil
}
