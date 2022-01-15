// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.19.3
// source: person_message.proto

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

type Person struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id    string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Shirt *Shirt `protobuf:"bytes,2,opt,name=shirt,proto3" json:"shirt,omitempty"`
	Chino *Chino `protobuf:"bytes,3,opt,name=chino,proto3" json:"chino,omitempty"`
}

func (x *Person) Reset() {
	*x = Person{}
	if protoimpl.UnsafeEnabled {
		mi := &file_person_message_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Person) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Person) ProtoMessage() {}

func (x *Person) ProtoReflect() protoreflect.Message {
	mi := &file_person_message_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Person.ProtoReflect.Descriptor instead.
func (*Person) Descriptor() ([]byte, []int) {
	return file_person_message_proto_rawDescGZIP(), []int{0}
}

func (x *Person) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Person) GetShirt() *Shirt {
	if x != nil {
		return x.Shirt
	}
	return nil
}

func (x *Person) GetChino() *Chino {
	if x != nil {
		return x.Chino
	}
	return nil
}

var File_person_message_proto protoreflect.FileDescriptor

var file_person_message_proto_rawDesc = []byte{
	0x0a, 0x14, 0x70, 0x65, 0x72, 0x73, 0x6f, 0x6e, 0x5f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x13, 0x73, 0x68, 0x69, 0x72, 0x74, 0x5f, 0x6d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x13, 0x63, 0x68, 0x69,
	0x6e, 0x6f, 0x5f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x22, 0x54, 0x0a, 0x06, 0x50, 0x65, 0x72, 0x73, 0x6f, 0x6e, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x1c, 0x0a, 0x05, 0x73, 0x68,
	0x69, 0x72, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x06, 0x2e, 0x53, 0x68, 0x69, 0x72,
	0x74, 0x52, 0x05, 0x73, 0x68, 0x69, 0x72, 0x74, 0x12, 0x1c, 0x0a, 0x05, 0x63, 0x68, 0x69, 0x6e,
	0x6f, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x06, 0x2e, 0x43, 0x68, 0x69, 0x6e, 0x6f, 0x52,
	0x05, 0x63, 0x68, 0x69, 0x6e, 0x6f, 0x42, 0x06, 0x5a, 0x04, 0x2e, 0x2f, 0x70, 0x62, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_person_message_proto_rawDescOnce sync.Once
	file_person_message_proto_rawDescData = file_person_message_proto_rawDesc
)

func file_person_message_proto_rawDescGZIP() []byte {
	file_person_message_proto_rawDescOnce.Do(func() {
		file_person_message_proto_rawDescData = protoimpl.X.CompressGZIP(file_person_message_proto_rawDescData)
	})
	return file_person_message_proto_rawDescData
}

var file_person_message_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_person_message_proto_goTypes = []interface{}{
	(*Person)(nil), // 0: Person
	(*Shirt)(nil),  // 1: Shirt
	(*Chino)(nil),  // 2: Chino
}
var file_person_message_proto_depIdxs = []int32{
	1, // 0: Person.shirt:type_name -> Shirt
	2, // 1: Person.chino:type_name -> Chino
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_person_message_proto_init() }
func file_person_message_proto_init() {
	if File_person_message_proto != nil {
		return
	}
	file_shirt_message_proto_init()
	file_chino_message_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_person_message_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Person); i {
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
			RawDescriptor: file_person_message_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_person_message_proto_goTypes,
		DependencyIndexes: file_person_message_proto_depIdxs,
		MessageInfos:      file_person_message_proto_msgTypes,
	}.Build()
	File_person_message_proto = out.File
	file_person_message_proto_rawDesc = nil
	file_person_message_proto_goTypes = nil
	file_person_message_proto_depIdxs = nil
}
