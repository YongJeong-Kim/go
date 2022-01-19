// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.19.3
// source: msg_service.proto

package pb

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

type SendUserRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Msg *Msg `protobuf:"bytes,1,opt,name=msg,proto3" json:"msg,omitempty"`
}

func (x *SendUserRequest) Reset() {
	*x = SendUserRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_msg_service_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SendUserRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SendUserRequest) ProtoMessage() {}

func (x *SendUserRequest) ProtoReflect() protoreflect.Message {
	mi := &file_msg_service_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SendUserRequest.ProtoReflect.Descriptor instead.
func (*SendUserRequest) Descriptor() ([]byte, []int) {
	return file_msg_service_proto_rawDescGZIP(), []int{0}
}

func (x *SendUserRequest) GetMsg() *Msg {
	if x != nil {
		return x.Msg
	}
	return nil
}

type SendUserResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id      string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	From    string `protobuf:"bytes,2,opt,name=from,proto3" json:"from,omitempty"`
	Content string `protobuf:"bytes,3,opt,name=content,proto3" json:"content,omitempty"`
}

func (x *SendUserResponse) Reset() {
	*x = SendUserResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_msg_service_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SendUserResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SendUserResponse) ProtoMessage() {}

func (x *SendUserResponse) ProtoReflect() protoreflect.Message {
	mi := &file_msg_service_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SendUserResponse.ProtoReflect.Descriptor instead.
func (*SendUserResponse) Descriptor() ([]byte, []int) {
	return file_msg_service_proto_rawDescGZIP(), []int{1}
}

func (x *SendUserResponse) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *SendUserResponse) GetFrom() string {
	if x != nil {
		return x.From
	}
	return ""
}

func (x *SendUserResponse) GetContent() string {
	if x != nil {
		return x.Content
	}
	return ""
}

var File_msg_service_proto protoreflect.FileDescriptor

var file_msg_service_proto_rawDesc = []byte{
	0x0a, 0x11, 0x6d, 0x73, 0x67, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x1a, 0x11, 0x6d, 0x73, 0x67, 0x5f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x29, 0x0a, 0x0f, 0x53, 0x65, 0x6e, 0x64, 0x55, 0x73,
	0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x16, 0x0a, 0x03, 0x6d, 0x73, 0x67,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x04, 0x2e, 0x4d, 0x73, 0x67, 0x52, 0x03, 0x6d, 0x73,
	0x67, 0x22, 0x50, 0x0a, 0x10, 0x53, 0x65, 0x6e, 0x64, 0x55, 0x73, 0x65, 0x72, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x66, 0x72, 0x6f, 0x6d, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x04, 0x66, 0x72, 0x6f, 0x6d, 0x12, 0x18, 0x0a, 0x07, 0x63, 0x6f, 0x6e,
	0x74, 0x65, 0x6e, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x63, 0x6f, 0x6e, 0x74,
	0x65, 0x6e, 0x74, 0x32, 0x45, 0x0a, 0x0a, 0x4d, 0x73, 0x67, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x12, 0x37, 0x0a, 0x0a, 0x53, 0x65, 0x6e, 0x64, 0x54, 0x6f, 0x55, 0x73, 0x65, 0x72, 0x12,
	0x10, 0x2e, 0x53, 0x65, 0x6e, 0x64, 0x55, 0x73, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x11, 0x2e, 0x53, 0x65, 0x6e, 0x64, 0x55, 0x73, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x28, 0x01, 0x30, 0x01, 0x42, 0x06, 0x5a, 0x04, 0x2e, 0x2f,
	0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_msg_service_proto_rawDescOnce sync.Once
	file_msg_service_proto_rawDescData = file_msg_service_proto_rawDesc
)

func file_msg_service_proto_rawDescGZIP() []byte {
	file_msg_service_proto_rawDescOnce.Do(func() {
		file_msg_service_proto_rawDescData = protoimpl.X.CompressGZIP(file_msg_service_proto_rawDescData)
	})
	return file_msg_service_proto_rawDescData
}

var file_msg_service_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_msg_service_proto_goTypes = []interface{}{
	(*SendUserRequest)(nil),  // 0: SendUserRequest
	(*SendUserResponse)(nil), // 1: SendUserResponse
	(*Msg)(nil),              // 2: Msg
}
var file_msg_service_proto_depIdxs = []int32{
	2, // 0: SendUserRequest.msg:type_name -> Msg
	0, // 1: MsgService.SendToUser:input_type -> SendUserRequest
	1, // 2: MsgService.SendToUser:output_type -> SendUserResponse
	2, // [2:3] is the sub-list for method output_type
	1, // [1:2] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_msg_service_proto_init() }
func file_msg_service_proto_init() {
	if File_msg_service_proto != nil {
		return
	}
	file_msg_message_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_msg_service_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SendUserRequest); i {
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
		file_msg_service_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SendUserResponse); i {
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
			RawDescriptor: file_msg_service_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_msg_service_proto_goTypes,
		DependencyIndexes: file_msg_service_proto_depIdxs,
		MessageInfos:      file_msg_service_proto_msgTypes,
	}.Build()
	File_msg_service_proto = out.File
	file_msg_service_proto_rawDesc = nil
	file_msg_service_proto_goTypes = nil
	file_msg_service_proto_depIdxs = nil
}
