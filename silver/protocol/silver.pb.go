// zinc balance protocol
// https://grpc.io/docs/languages/go/quickstart/

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v4.25.1
// source: silver.proto

package protocol

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

// repository basic type
type Repository struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// repo id
	Id int64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	// relative path not disk location
	RelativePath string `protobuf:"bytes,2,opt,name=relative_path,json=relativePath,proto3" json:"relative_path,omitempty"`
}

func (x *Repository) Reset() {
	*x = Repository{}
	if protoimpl.UnsafeEnabled {
		mi := &file_silver_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Repository) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Repository) ProtoMessage() {}

func (x *Repository) ProtoReflect() protoreflect.Message {
	mi := &file_silver_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Repository.ProtoReflect.Descriptor instead.
func (*Repository) Descriptor() ([]byte, []int) {
	return file_silver_proto_rawDescGZIP(), []int{0}
}

func (x *Repository) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Repository) GetRelativePath() string {
	if x != nil {
		return x.RelativePath
	}
	return ""
}

// Git over SSH fetch/clone request
type UploadPackRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Repo *Repository `protobuf:"bytes,1,opt,name=repo,proto3" json:"repo,omitempty"`
	Uid  int64       `protobuf:"varint,2,opt,name=uid,proto3" json:"uid,omitempty"`
	// git-upload-pack stdin
	Stdin []byte `protobuf:"bytes,3,opt,name=stdin,proto3" json:"stdin,omitempty"`
	// eg: version=2
	Protocol string `protobuf:"bytes,4,opt,name=protocol,proto3" json:"protocol,omitempty"`
}

func (x *UploadPackRequest) Reset() {
	*x = UploadPackRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_silver_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UploadPackRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UploadPackRequest) ProtoMessage() {}

func (x *UploadPackRequest) ProtoReflect() protoreflect.Message {
	mi := &file_silver_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UploadPackRequest.ProtoReflect.Descriptor instead.
func (*UploadPackRequest) Descriptor() ([]byte, []int) {
	return file_silver_proto_rawDescGZIP(), []int{1}
}

func (x *UploadPackRequest) GetRepo() *Repository {
	if x != nil {
		return x.Repo
	}
	return nil
}

func (x *UploadPackRequest) GetUid() int64 {
	if x != nil {
		return x.Uid
	}
	return 0
}

func (x *UploadPackRequest) GetStdin() []byte {
	if x != nil {
		return x.Stdin
	}
	return nil
}

func (x *UploadPackRequest) GetProtocol() string {
	if x != nil {
		return x.Protocol
	}
	return ""
}

// Git over SSH fetch/clone response
type UploadPackResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// git-upload-pack stdout
	Stdout []byte `protobuf:"bytes,1,opt,name=stdout,proto3" json:"stdout,omitempty"`
	// git-upload-pack stderr
	Stderr []byte `protobuf:"bytes,2,opt,name=stderr,proto3" json:"stderr,omitempty"`
	// git-upload-pack exit code
	ExitCode int32 `protobuf:"varint,3,opt,name=exit_code,json=exitCode,proto3" json:"exit_code,omitempty"`
}

func (x *UploadPackResponse) Reset() {
	*x = UploadPackResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_silver_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UploadPackResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UploadPackResponse) ProtoMessage() {}

func (x *UploadPackResponse) ProtoReflect() protoreflect.Message {
	mi := &file_silver_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UploadPackResponse.ProtoReflect.Descriptor instead.
func (*UploadPackResponse) Descriptor() ([]byte, []int) {
	return file_silver_proto_rawDescGZIP(), []int{2}
}

func (x *UploadPackResponse) GetStdout() []byte {
	if x != nil {
		return x.Stdout
	}
	return nil
}

func (x *UploadPackResponse) GetStderr() []byte {
	if x != nil {
		return x.Stderr
	}
	return nil
}

func (x *UploadPackResponse) GetExitCode() int32 {
	if x != nil {
		return x.ExitCode
	}
	return 0
}

// Git over SSH push request
type ReceivePackRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Repo *Repository `protobuf:"bytes,1,opt,name=repo,proto3" json:"repo,omitempty"`
	Uid  int64       `protobuf:"varint,2,opt,name=uid,proto3" json:"uid,omitempty"`
	// git-receive-pack stdin
	Stdin []byte `protobuf:"bytes,3,opt,name=stdin,proto3" json:"stdin,omitempty"`
	// eg: version=2
	Protocol string `protobuf:"bytes,4,opt,name=protocol,proto3" json:"protocol,omitempty"`
}

func (x *ReceivePackRequest) Reset() {
	*x = ReceivePackRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_silver_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReceivePackRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReceivePackRequest) ProtoMessage() {}

func (x *ReceivePackRequest) ProtoReflect() protoreflect.Message {
	mi := &file_silver_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReceivePackRequest.ProtoReflect.Descriptor instead.
func (*ReceivePackRequest) Descriptor() ([]byte, []int) {
	return file_silver_proto_rawDescGZIP(), []int{3}
}

func (x *ReceivePackRequest) GetRepo() *Repository {
	if x != nil {
		return x.Repo
	}
	return nil
}

func (x *ReceivePackRequest) GetUid() int64 {
	if x != nil {
		return x.Uid
	}
	return 0
}

func (x *ReceivePackRequest) GetStdin() []byte {
	if x != nil {
		return x.Stdin
	}
	return nil
}

func (x *ReceivePackRequest) GetProtocol() string {
	if x != nil {
		return x.Protocol
	}
	return ""
}

// Git over SSH push request
type ReceivePackResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// git-receive-pack stdout
	Stdout []byte `protobuf:"bytes,1,opt,name=stdout,proto3" json:"stdout,omitempty"`
	// git-receive-pack stderr
	Stderr []byte `protobuf:"bytes,2,opt,name=stderr,proto3" json:"stderr,omitempty"`
	// git-receive-pack exit code
	ExitCode int32 `protobuf:"varint,3,opt,name=exit_code,json=exitCode,proto3" json:"exit_code,omitempty"`
}

func (x *ReceivePackResponse) Reset() {
	*x = ReceivePackResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_silver_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReceivePackResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReceivePackResponse) ProtoMessage() {}

func (x *ReceivePackResponse) ProtoReflect() protoreflect.Message {
	mi := &file_silver_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReceivePackResponse.ProtoReflect.Descriptor instead.
func (*ReceivePackResponse) Descriptor() ([]byte, []int) {
	return file_silver_proto_rawDescGZIP(), []int{4}
}

func (x *ReceivePackResponse) GetStdout() []byte {
	if x != nil {
		return x.Stdout
	}
	return nil
}

func (x *ReceivePackResponse) GetStderr() []byte {
	if x != nil {
		return x.Stderr
	}
	return nil
}

func (x *ReceivePackResponse) GetExitCode() int32 {
	if x != nil {
		return x.ExitCode
	}
	return 0
}

// Git over HTTP get refs request
type InfoRefsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Repo        *Repository `protobuf:"bytes,1,opt,name=repo,proto3" json:"repo,omitempty"`
	Uid         int64       `protobuf:"varint,2,opt,name=uid,proto3" json:"uid,omitempty"`
	ServiceName string      `protobuf:"bytes,3,opt,name=service_name,json=serviceName,proto3" json:"service_name,omitempty"`
	// eg: version=2
	Protocol string `protobuf:"bytes,4,opt,name=protocol,proto3" json:"protocol,omitempty"`
}

func (x *InfoRefsRequest) Reset() {
	*x = InfoRefsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_silver_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *InfoRefsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InfoRefsRequest) ProtoMessage() {}

func (x *InfoRefsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_silver_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InfoRefsRequest.ProtoReflect.Descriptor instead.
func (*InfoRefsRequest) Descriptor() ([]byte, []int) {
	return file_silver_proto_rawDescGZIP(), []int{5}
}

func (x *InfoRefsRequest) GetRepo() *Repository {
	if x != nil {
		return x.Repo
	}
	return nil
}

func (x *InfoRefsRequest) GetUid() int64 {
	if x != nil {
		return x.Uid
	}
	return 0
}

func (x *InfoRefsRequest) GetServiceName() string {
	if x != nil {
		return x.ServiceName
	}
	return ""
}

func (x *InfoRefsRequest) GetProtocol() string {
	if x != nil {
		return x.Protocol
	}
	return ""
}

type UploadArchiveRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// 'repository' must be present in the first message.
	Repo *Repository `protobuf:"bytes,1,opt,name=repo,proto3" json:"repo,omitempty"`
	// A chunk of raw data to be copied to 'git upload-archive' standard input
	Stdin []byte `protobuf:"bytes,2,opt,name=stdin,proto3" json:"stdin,omitempty"`
}

func (x *UploadArchiveRequest) Reset() {
	*x = UploadArchiveRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_silver_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UploadArchiveRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UploadArchiveRequest) ProtoMessage() {}

func (x *UploadArchiveRequest) ProtoReflect() protoreflect.Message {
	mi := &file_silver_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UploadArchiveRequest.ProtoReflect.Descriptor instead.
func (*UploadArchiveRequest) Descriptor() ([]byte, []int) {
	return file_silver_proto_rawDescGZIP(), []int{6}
}

func (x *UploadArchiveRequest) GetRepo() *Repository {
	if x != nil {
		return x.Repo
	}
	return nil
}

func (x *UploadArchiveRequest) GetStdin() []byte {
	if x != nil {
		return x.Stdin
	}
	return nil
}

type UploadArchiveResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// A chunk of raw data from 'git upload-archive' standard output
	Stdout []byte `protobuf:"bytes,1,opt,name=stdout,proto3" json:"stdout,omitempty"`
	// A chunk of raw data from 'git upload-archive' standard error
	Stderr []byte `protobuf:"bytes,2,opt,name=stderr,proto3" json:"stderr,omitempty"`
	// This value will only be set on the last message
	ExitCode int32 `protobuf:"varint,3,opt,name=exit_code,json=exitCode,proto3" json:"exit_code,omitempty"`
}

func (x *UploadArchiveResponse) Reset() {
	*x = UploadArchiveResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_silver_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UploadArchiveResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UploadArchiveResponse) ProtoMessage() {}

func (x *UploadArchiveResponse) ProtoReflect() protoreflect.Message {
	mi := &file_silver_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UploadArchiveResponse.ProtoReflect.Descriptor instead.
func (*UploadArchiveResponse) Descriptor() ([]byte, []int) {
	return file_silver_proto_rawDescGZIP(), []int{7}
}

func (x *UploadArchiveResponse) GetStdout() []byte {
	if x != nil {
		return x.Stdout
	}
	return nil
}

func (x *UploadArchiveResponse) GetStderr() []byte {
	if x != nil {
		return x.Stderr
	}
	return nil
}

func (x *UploadArchiveResponse) GetExitCode() int32 {
	if x != nil {
		return x.ExitCode
	}
	return 0
}

// Git over HTTP get refs response
type InfoRefsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// git-upload-pack stdout
	// git-receive-pack stdout
	Stdout []byte `protobuf:"bytes,1,opt,name=stdout,proto3" json:"stdout,omitempty"`
}

func (x *InfoRefsResponse) Reset() {
	*x = InfoRefsResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_silver_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *InfoRefsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InfoRefsResponse) ProtoMessage() {}

func (x *InfoRefsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_silver_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InfoRefsResponse.ProtoReflect.Descriptor instead.
func (*InfoRefsResponse) Descriptor() ([]byte, []int) {
	return file_silver_proto_rawDescGZIP(), []int{8}
}

func (x *InfoRefsResponse) GetStdout() []byte {
	if x != nil {
		return x.Stdout
	}
	return nil
}

// Git Over HTTP fetch/clone request (POST request body)
type PostUploadPackRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Repo *Repository `protobuf:"bytes,1,opt,name=repo,proto3" json:"repo,omitempty"`
	// git-upload-pack stdin
	Stdin []byte `protobuf:"bytes,2,opt,name=stdin,proto3" json:"stdin,omitempty"`
	// eg: version=2
	Protocol string `protobuf:"bytes,3,opt,name=protocol,proto3" json:"protocol,omitempty"`
}

func (x *PostUploadPackRequest) Reset() {
	*x = PostUploadPackRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_silver_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PostUploadPackRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PostUploadPackRequest) ProtoMessage() {}

func (x *PostUploadPackRequest) ProtoReflect() protoreflect.Message {
	mi := &file_silver_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PostUploadPackRequest.ProtoReflect.Descriptor instead.
func (*PostUploadPackRequest) Descriptor() ([]byte, []int) {
	return file_silver_proto_rawDescGZIP(), []int{9}
}

func (x *PostUploadPackRequest) GetRepo() *Repository {
	if x != nil {
		return x.Repo
	}
	return nil
}

func (x *PostUploadPackRequest) GetStdin() []byte {
	if x != nil {
		return x.Stdin
	}
	return nil
}

func (x *PostUploadPackRequest) GetProtocol() string {
	if x != nil {
		return x.Protocol
	}
	return ""
}

// Git Over HTTP fetch/clone response (POST response body)
type PostUploadPackResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// git-upload-pack stdout
	Stdout []byte `protobuf:"bytes,1,opt,name=stdout,proto3" json:"stdout,omitempty"`
}

func (x *PostUploadPackResponse) Reset() {
	*x = PostUploadPackResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_silver_proto_msgTypes[10]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PostUploadPackResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PostUploadPackResponse) ProtoMessage() {}

func (x *PostUploadPackResponse) ProtoReflect() protoreflect.Message {
	mi := &file_silver_proto_msgTypes[10]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PostUploadPackResponse.ProtoReflect.Descriptor instead.
func (*PostUploadPackResponse) Descriptor() ([]byte, []int) {
	return file_silver_proto_rawDescGZIP(), []int{10}
}

func (x *PostUploadPackResponse) GetStdout() []byte {
	if x != nil {
		return x.Stdout
	}
	return nil
}

// Git Over HTTP push request (POST request body)
type PostReceivePackRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Repo *Repository `protobuf:"bytes,1,opt,name=repo,proto3" json:"repo,omitempty"`
	Uid  int64       `protobuf:"varint,2,opt,name=uid,proto3" json:"uid,omitempty"`
	// git-receive-pack stdin
	Stdin []byte `protobuf:"bytes,3,opt,name=stdin,proto3" json:"stdin,omitempty"`
	// eg: version=2
	Protocol string `protobuf:"bytes,4,opt,name=protocol,proto3" json:"protocol,omitempty"`
}

func (x *PostReceivePackRequest) Reset() {
	*x = PostReceivePackRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_silver_proto_msgTypes[11]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PostReceivePackRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PostReceivePackRequest) ProtoMessage() {}

func (x *PostReceivePackRequest) ProtoReflect() protoreflect.Message {
	mi := &file_silver_proto_msgTypes[11]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PostReceivePackRequest.ProtoReflect.Descriptor instead.
func (*PostReceivePackRequest) Descriptor() ([]byte, []int) {
	return file_silver_proto_rawDescGZIP(), []int{11}
}

func (x *PostReceivePackRequest) GetRepo() *Repository {
	if x != nil {
		return x.Repo
	}
	return nil
}

func (x *PostReceivePackRequest) GetUid() int64 {
	if x != nil {
		return x.Uid
	}
	return 0
}

func (x *PostReceivePackRequest) GetStdin() []byte {
	if x != nil {
		return x.Stdin
	}
	return nil
}

func (x *PostReceivePackRequest) GetProtocol() string {
	if x != nil {
		return x.Protocol
	}
	return ""
}

// Git Over HTTP push response (POST response body)
type PostReceivePackResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// git-receive-pack stdout
	Stdout []byte `protobuf:"bytes,1,opt,name=stdout,proto3" json:"stdout,omitempty"`
}

func (x *PostReceivePackResponse) Reset() {
	*x = PostReceivePackResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_silver_proto_msgTypes[12]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PostReceivePackResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PostReceivePackResponse) ProtoMessage() {}

func (x *PostReceivePackResponse) ProtoReflect() protoreflect.Message {
	mi := &file_silver_proto_msgTypes[12]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PostReceivePackResponse.ProtoReflect.Descriptor instead.
func (*PostReceivePackResponse) Descriptor() ([]byte, []int) {
	return file_silver_proto_rawDescGZIP(), []int{12}
}

func (x *PostReceivePackResponse) GetStdout() []byte {
	if x != nil {
		return x.Stdout
	}
	return nil
}

var File_silver_proto protoreflect.FileDescriptor

var file_silver_proto_rawDesc = []byte{
	0x0a, 0x0c, 0x73, 0x69, 0x6c, 0x76, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x41,
	0x0a, 0x0a, 0x52, 0x65, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x6f, 0x72, 0x79, 0x12, 0x0e, 0x0a, 0x02,
	0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x69, 0x64, 0x12, 0x23, 0x0a, 0x0d,
	0x72, 0x65, 0x6c, 0x61, 0x74, 0x69, 0x76, 0x65, 0x5f, 0x70, 0x61, 0x74, 0x68, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0c, 0x72, 0x65, 0x6c, 0x61, 0x74, 0x69, 0x76, 0x65, 0x50, 0x61, 0x74,
	0x68, 0x22, 0x78, 0x0a, 0x11, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x50, 0x61, 0x63, 0x6b, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1f, 0x0a, 0x04, 0x72, 0x65, 0x70, 0x6f, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x0b, 0x2e, 0x52, 0x65, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x6f, 0x72,
	0x79, 0x52, 0x04, 0x72, 0x65, 0x70, 0x6f, 0x12, 0x10, 0x0a, 0x03, 0x75, 0x69, 0x64, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x03, 0x52, 0x03, 0x75, 0x69, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x73, 0x74, 0x64,
	0x69, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x05, 0x73, 0x74, 0x64, 0x69, 0x6e, 0x12,
	0x1a, 0x0a, 0x08, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x08, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x22, 0x61, 0x0a, 0x12, 0x55,
	0x70, 0x6c, 0x6f, 0x61, 0x64, 0x50, 0x61, 0x63, 0x6b, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x74, 0x64, 0x6f, 0x75, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0c, 0x52, 0x06, 0x73, 0x74, 0x64, 0x6f, 0x75, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x74, 0x64,
	0x65, 0x72, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x06, 0x73, 0x74, 0x64, 0x65, 0x72,
	0x72, 0x12, 0x1b, 0x0a, 0x09, 0x65, 0x78, 0x69, 0x74, 0x5f, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x05, 0x52, 0x08, 0x65, 0x78, 0x69, 0x74, 0x43, 0x6f, 0x64, 0x65, 0x22, 0x79,
	0x0a, 0x12, 0x52, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x50, 0x61, 0x63, 0x6b, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x1f, 0x0a, 0x04, 0x72, 0x65, 0x70, 0x6f, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x0b, 0x2e, 0x52, 0x65, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x6f, 0x72, 0x79, 0x52,
	0x04, 0x72, 0x65, 0x70, 0x6f, 0x12, 0x10, 0x0a, 0x03, 0x75, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x03, 0x52, 0x03, 0x75, 0x69, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x73, 0x74, 0x64, 0x69, 0x6e,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x05, 0x73, 0x74, 0x64, 0x69, 0x6e, 0x12, 0x1a, 0x0a,
	0x08, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x08, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x22, 0x62, 0x0a, 0x13, 0x52, 0x65, 0x63,
	0x65, 0x69, 0x76, 0x65, 0x50, 0x61, 0x63, 0x6b, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x16, 0x0a, 0x06, 0x73, 0x74, 0x64, 0x6f, 0x75, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c,
	0x52, 0x06, 0x73, 0x74, 0x64, 0x6f, 0x75, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x74, 0x64, 0x65,
	0x72, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x06, 0x73, 0x74, 0x64, 0x65, 0x72, 0x72,
	0x12, 0x1b, 0x0a, 0x09, 0x65, 0x78, 0x69, 0x74, 0x5f, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x05, 0x52, 0x08, 0x65, 0x78, 0x69, 0x74, 0x43, 0x6f, 0x64, 0x65, 0x22, 0x83, 0x01,
	0x0a, 0x0f, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x65, 0x66, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x1f, 0x0a, 0x04, 0x72, 0x65, 0x70, 0x6f, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x0b, 0x2e, 0x52, 0x65, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x6f, 0x72, 0x79, 0x52, 0x04, 0x72, 0x65,
	0x70, 0x6f, 0x12, 0x10, 0x0a, 0x03, 0x75, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x03, 0x75, 0x69, 0x64, 0x12, 0x21, 0x0a, 0x0c, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x5f,
	0x6e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x73, 0x65, 0x72, 0x76,
	0x69, 0x63, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x63, 0x6f, 0x6c, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x63, 0x6f, 0x6c, 0x22, 0x4d, 0x0a, 0x14, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x41, 0x72, 0x63,
	0x68, 0x69, 0x76, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1f, 0x0a, 0x04, 0x72,
	0x65, 0x70, 0x6f, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0b, 0x2e, 0x52, 0x65, 0x70, 0x6f,
	0x73, 0x69, 0x74, 0x6f, 0x72, 0x79, 0x52, 0x04, 0x72, 0x65, 0x70, 0x6f, 0x12, 0x14, 0x0a, 0x05,
	0x73, 0x74, 0x64, 0x69, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x05, 0x73, 0x74, 0x64,
	0x69, 0x6e, 0x22, 0x64, 0x0a, 0x15, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x41, 0x72, 0x63, 0x68,
	0x69, 0x76, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x73,
	0x74, 0x64, 0x6f, 0x75, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x06, 0x73, 0x74, 0x64,
	0x6f, 0x75, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x74, 0x64, 0x65, 0x72, 0x72, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x0c, 0x52, 0x06, 0x73, 0x74, 0x64, 0x65, 0x72, 0x72, 0x12, 0x1b, 0x0a, 0x09, 0x65,
	0x78, 0x69, 0x74, 0x5f, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x08,
	0x65, 0x78, 0x69, 0x74, 0x43, 0x6f, 0x64, 0x65, 0x22, 0x2a, 0x0a, 0x10, 0x49, 0x6e, 0x66, 0x6f,
	0x52, 0x65, 0x66, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x16, 0x0a, 0x06,
	0x73, 0x74, 0x64, 0x6f, 0x75, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x06, 0x73, 0x74,
	0x64, 0x6f, 0x75, 0x74, 0x22, 0x6a, 0x0a, 0x15, 0x50, 0x6f, 0x73, 0x74, 0x55, 0x70, 0x6c, 0x6f,
	0x61, 0x64, 0x50, 0x61, 0x63, 0x6b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1f, 0x0a,
	0x04, 0x72, 0x65, 0x70, 0x6f, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0b, 0x2e, 0x52, 0x65,
	0x70, 0x6f, 0x73, 0x69, 0x74, 0x6f, 0x72, 0x79, 0x52, 0x04, 0x72, 0x65, 0x70, 0x6f, 0x12, 0x14,
	0x0a, 0x05, 0x73, 0x74, 0x64, 0x69, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x05, 0x73,
	0x74, 0x64, 0x69, 0x6e, 0x12, 0x1a, 0x0a, 0x08, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c,
	0x22, 0x30, 0x0a, 0x16, 0x50, 0x6f, 0x73, 0x74, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x50, 0x61,
	0x63, 0x6b, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x74,
	0x64, 0x6f, 0x75, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x06, 0x73, 0x74, 0x64, 0x6f,
	0x75, 0x74, 0x22, 0x7d, 0x0a, 0x16, 0x50, 0x6f, 0x73, 0x74, 0x52, 0x65, 0x63, 0x65, 0x69, 0x76,
	0x65, 0x50, 0x61, 0x63, 0x6b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1f, 0x0a, 0x04,
	0x72, 0x65, 0x70, 0x6f, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0b, 0x2e, 0x52, 0x65, 0x70,
	0x6f, 0x73, 0x69, 0x74, 0x6f, 0x72, 0x79, 0x52, 0x04, 0x72, 0x65, 0x70, 0x6f, 0x12, 0x10, 0x0a,
	0x03, 0x75, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x03, 0x75, 0x69, 0x64, 0x12,
	0x14, 0x0a, 0x05, 0x73, 0x74, 0x64, 0x69, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x05,
	0x73, 0x74, 0x64, 0x69, 0x6e, 0x12, 0x1a, 0x0a, 0x08, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f,
	0x6c, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f,
	0x6c, 0x22, 0x31, 0x0a, 0x17, 0x50, 0x6f, 0x73, 0x74, 0x52, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65,
	0x50, 0x61, 0x63, 0x6b, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x16, 0x0a, 0x06,
	0x73, 0x74, 0x64, 0x6f, 0x75, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x06, 0x73, 0x74,
	0x64, 0x6f, 0x75, 0x74, 0x32, 0xd1, 0x03, 0x0a, 0x06, 0x53, 0x69, 0x6c, 0x76, 0x65, 0x72, 0x12,
	0x39, 0x0a, 0x0a, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x50, 0x61, 0x63, 0x6b, 0x12, 0x12, 0x2e,
	0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x50, 0x61, 0x63, 0x6b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x13, 0x2e, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x50, 0x61, 0x63, 0x6b, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x28, 0x01, 0x30, 0x01, 0x12, 0x3c, 0x0a, 0x0b, 0x52, 0x65,
	0x63, 0x65, 0x69, 0x76, 0x65, 0x50, 0x61, 0x63, 0x6b, 0x12, 0x13, 0x2e, 0x52, 0x65, 0x63, 0x65,
	0x69, 0x76, 0x65, 0x50, 0x61, 0x63, 0x6b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x14,
	0x2e, 0x52, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x50, 0x61, 0x63, 0x6b, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x28, 0x01, 0x30, 0x01, 0x12, 0x42, 0x0a, 0x0d, 0x55, 0x70, 0x6c, 0x6f,
	0x61, 0x64, 0x41, 0x72, 0x63, 0x68, 0x69, 0x76, 0x65, 0x12, 0x15, 0x2e, 0x55, 0x70, 0x6c, 0x6f,
	0x61, 0x64, 0x41, 0x72, 0x63, 0x68, 0x69, 0x76, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x16, 0x2e, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x41, 0x72, 0x63, 0x68, 0x69, 0x76, 0x65,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x28, 0x01, 0x30, 0x01, 0x12, 0x3b, 0x0a, 0x12,
	0x49, 0x6e, 0x66, 0x6f, 0x52, 0x65, 0x66, 0x73, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x50, 0x61,
	0x63, 0x6b, 0x12, 0x10, 0x2e, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x65, 0x66, 0x73, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x11, 0x2e, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x65, 0x66, 0x73, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x30, 0x01, 0x12, 0x3c, 0x0a, 0x13, 0x49, 0x6e, 0x66,
	0x6f, 0x52, 0x65, 0x66, 0x73, 0x52, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x50, 0x61, 0x63, 0x6b,
	0x12, 0x10, 0x2e, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x65, 0x66, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x11, 0x2e, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x65, 0x66, 0x73, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x30, 0x01, 0x12, 0x45, 0x0a, 0x0e, 0x50, 0x6f, 0x73, 0x74, 0x55,
	0x70, 0x6c, 0x6f, 0x61, 0x64, 0x50, 0x61, 0x63, 0x6b, 0x12, 0x16, 0x2e, 0x50, 0x6f, 0x73, 0x74,
	0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x50, 0x61, 0x63, 0x6b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x17, 0x2e, 0x50, 0x6f, 0x73, 0x74, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x50, 0x61,
	0x63, 0x6b, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x28, 0x01, 0x30, 0x01, 0x12, 0x48,
	0x0a, 0x0f, 0x50, 0x6f, 0x73, 0x74, 0x52, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x50, 0x61, 0x63,
	0x6b, 0x12, 0x17, 0x2e, 0x50, 0x6f, 0x73, 0x74, 0x52, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x50,
	0x61, 0x63, 0x6b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x18, 0x2e, 0x50, 0x6f, 0x73,
	0x74, 0x52, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x50, 0x61, 0x63, 0x6b, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x28, 0x01, 0x30, 0x01, 0x42, 0x29, 0x5a, 0x27, 0x67, 0x69, 0x74, 0x68,
	0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x7a, 0x69, 0x6e, 0x63, 0x69, 0x75, 0x6d, 0x2f, 0x7a,
	0x69, 0x6e, 0x63, 0x2f, 0x73, 0x6c, 0x69, 0x76, 0x65, 0x72, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x63, 0x6f, 0x6c, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_silver_proto_rawDescOnce sync.Once
	file_silver_proto_rawDescData = file_silver_proto_rawDesc
)

func file_silver_proto_rawDescGZIP() []byte {
	file_silver_proto_rawDescOnce.Do(func() {
		file_silver_proto_rawDescData = protoimpl.X.CompressGZIP(file_silver_proto_rawDescData)
	})
	return file_silver_proto_rawDescData
}

var file_silver_proto_msgTypes = make([]protoimpl.MessageInfo, 13)
var file_silver_proto_goTypes = []interface{}{
	(*Repository)(nil),              // 0: Repository
	(*UploadPackRequest)(nil),       // 1: UploadPackRequest
	(*UploadPackResponse)(nil),      // 2: UploadPackResponse
	(*ReceivePackRequest)(nil),      // 3: ReceivePackRequest
	(*ReceivePackResponse)(nil),     // 4: ReceivePackResponse
	(*InfoRefsRequest)(nil),         // 5: InfoRefsRequest
	(*UploadArchiveRequest)(nil),    // 6: UploadArchiveRequest
	(*UploadArchiveResponse)(nil),   // 7: UploadArchiveResponse
	(*InfoRefsResponse)(nil),        // 8: InfoRefsResponse
	(*PostUploadPackRequest)(nil),   // 9: PostUploadPackRequest
	(*PostUploadPackResponse)(nil),  // 10: PostUploadPackResponse
	(*PostReceivePackRequest)(nil),  // 11: PostReceivePackRequest
	(*PostReceivePackResponse)(nil), // 12: PostReceivePackResponse
}
var file_silver_proto_depIdxs = []int32{
	0,  // 0: UploadPackRequest.repo:type_name -> Repository
	0,  // 1: ReceivePackRequest.repo:type_name -> Repository
	0,  // 2: InfoRefsRequest.repo:type_name -> Repository
	0,  // 3: UploadArchiveRequest.repo:type_name -> Repository
	0,  // 4: PostUploadPackRequest.repo:type_name -> Repository
	0,  // 5: PostReceivePackRequest.repo:type_name -> Repository
	1,  // 6: Silver.UploadPack:input_type -> UploadPackRequest
	3,  // 7: Silver.ReceivePack:input_type -> ReceivePackRequest
	6,  // 8: Silver.UploadArchive:input_type -> UploadArchiveRequest
	5,  // 9: Silver.InfoRefsUploadPack:input_type -> InfoRefsRequest
	5,  // 10: Silver.InfoRefsReceivePack:input_type -> InfoRefsRequest
	9,  // 11: Silver.PostUploadPack:input_type -> PostUploadPackRequest
	11, // 12: Silver.PostReceivePack:input_type -> PostReceivePackRequest
	2,  // 13: Silver.UploadPack:output_type -> UploadPackResponse
	4,  // 14: Silver.ReceivePack:output_type -> ReceivePackResponse
	7,  // 15: Silver.UploadArchive:output_type -> UploadArchiveResponse
	8,  // 16: Silver.InfoRefsUploadPack:output_type -> InfoRefsResponse
	8,  // 17: Silver.InfoRefsReceivePack:output_type -> InfoRefsResponse
	10, // 18: Silver.PostUploadPack:output_type -> PostUploadPackResponse
	12, // 19: Silver.PostReceivePack:output_type -> PostReceivePackResponse
	13, // [13:20] is the sub-list for method output_type
	6,  // [6:13] is the sub-list for method input_type
	6,  // [6:6] is the sub-list for extension type_name
	6,  // [6:6] is the sub-list for extension extendee
	0,  // [0:6] is the sub-list for field type_name
}

func init() { file_silver_proto_init() }
func file_silver_proto_init() {
	if File_silver_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_silver_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Repository); i {
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
		file_silver_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UploadPackRequest); i {
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
		file_silver_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UploadPackResponse); i {
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
		file_silver_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ReceivePackRequest); i {
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
		file_silver_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ReceivePackResponse); i {
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
		file_silver_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*InfoRefsRequest); i {
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
		file_silver_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UploadArchiveRequest); i {
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
		file_silver_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UploadArchiveResponse); i {
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
		file_silver_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*InfoRefsResponse); i {
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
		file_silver_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PostUploadPackRequest); i {
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
		file_silver_proto_msgTypes[10].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PostUploadPackResponse); i {
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
		file_silver_proto_msgTypes[11].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PostReceivePackRequest); i {
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
		file_silver_proto_msgTypes[12].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PostReceivePackResponse); i {
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
			RawDescriptor: file_silver_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   13,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_silver_proto_goTypes,
		DependencyIndexes: file_silver_proto_depIdxs,
		MessageInfos:      file_silver_proto_msgTypes,
	}.Build()
	File_silver_proto = out.File
	file_silver_proto_rawDesc = nil
	file_silver_proto_goTypes = nil
	file_silver_proto_depIdxs = nil
}
