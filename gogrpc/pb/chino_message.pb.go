// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.19.3
// source: chino_message.proto

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

type Chino struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Brand string `protobuf:"bytes,1,opt,name=brand,proto3" json:"brand,omitempty"`
	Name  string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Color Color  `protobuf:"varint,3,opt,name=color,proto3,enum=Color" json:"color,omitempty"`
}

func (x *Chino) Reset() {
	*x = Chino{}
	if protoimpl.UnsafeEnabled {
		mi := &file_chino_message_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Chino) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Chino) ProtoMessage() {}

func (x *Chino) ProtoReflect() protoreflect.Message {
	mi := &file_chino_message_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Chino.ProtoReflect.Descriptor instead.
func (*Chino) Descriptor() ([]byte, []int) {
	return file_chino_message_proto_rawDescGZIP(), []int{0}
}

func (x *Chino) GetBrand() string {
	if x != nil {
		return x.Brand
	}
	return ""
}

func (x *Chino) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Chino) GetColor() Color {
	if x != nil {
		return x.Color
	}
	return Color_RED
}

var File_chino_message_proto protoreflect.FileDescriptor

var file_chino_message_proto_rawDesc = []byte{
	0x0a, 0x13, 0x63, 0x68, 0x69, 0x6e, 0x6f, 0x5f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x10, 0x63, 0x6f, 0x6c, 0x6f, 0x72, 0x5f, 0x65, 0x6e, 0x75,
	0x6d, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x4f, 0x0a, 0x05, 0x43, 0x68, 0x69, 0x6e, 0x6f,
	0x12, 0x14, 0x0a, 0x05, 0x62, 0x72, 0x61, 0x6e, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x05, 0x62, 0x72, 0x61, 0x6e, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x1c, 0x0a, 0x05, 0x63, 0x6f,
	0x6c, 0x6f, 0x72, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x06, 0x2e, 0x43, 0x6f, 0x6c, 0x6f,
	0x72, 0x52, 0x05, 0x63, 0x6f, 0x6c, 0x6f, 0x72, 0x42, 0x06, 0x5a, 0x04, 0x2e, 0x2f, 0x70, 0x62,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_chino_message_proto_rawDescOnce sync.Once
	file_chino_message_proto_rawDescData = file_chino_message_proto_rawDesc
)

func file_chino_message_proto_rawDescGZIP() []byte {
	file_chino_message_proto_rawDescOnce.Do(func() {
		file_chino_message_proto_rawDescData = protoimpl.X.CompressGZIP(file_chino_message_proto_rawDescData)
	})
	return file_chino_message_proto_rawDescData
}

var file_chino_message_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_chino_message_proto_goTypes = []interface{}{
	(*Chino)(nil), // 0: Chino
	(Color)(0),    // 1: Color
}
var file_chino_message_proto_depIdxs = []int32{
	1, // 0: Chino.color:type_name -> Color
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_chino_message_proto_init() }
func file_chino_message_proto_init() {
	if File_chino_message_proto != nil {
		return
	}
	file_color_enum_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_chino_message_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Chino); i {
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
			RawDescriptor: file_chino_message_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_chino_message_proto_goTypes,
		DependencyIndexes: file_chino_message_proto_depIdxs,
		MessageInfos:      file_chino_message_proto_msgTypes,
	}.Build()
	File_chino_message_proto = out.File
	file_chino_message_proto_rawDesc = nil
	file_chino_message_proto_goTypes = nil
	file_chino_message_proto_depIdxs = nil
}
