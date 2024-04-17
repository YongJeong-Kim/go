// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        (unknown)
// source: shop/v1/shop_message.proto

package shopv1

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

type Shop struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ShopId string `protobuf:"bytes,1,opt,name=shop_id,json=shopId,proto3" json:"shop_id,omitempty"`
}

func (x *Shop) Reset() {
	*x = Shop{}
	if protoimpl.UnsafeEnabled {
		mi := &file_shop_v1_shop_message_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Shop) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Shop) ProtoMessage() {}

func (x *Shop) ProtoReflect() protoreflect.Message {
	mi := &file_shop_v1_shop_message_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Shop.ProtoReflect.Descriptor instead.
func (*Shop) Descriptor() ([]byte, []int) {
	return file_shop_v1_shop_message_proto_rawDescGZIP(), []int{0}
}

func (x *Shop) GetShopId() string {
	if x != nil {
		return x.ShopId
	}
	return ""
}

type CreateShopRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Shop *Shop `protobuf:"bytes,1,opt,name=shop,proto3" json:"shop,omitempty"`
}

func (x *CreateShopRequest) Reset() {
	*x = CreateShopRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_shop_v1_shop_message_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateShopRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateShopRequest) ProtoMessage() {}

func (x *CreateShopRequest) ProtoReflect() protoreflect.Message {
	mi := &file_shop_v1_shop_message_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateShopRequest.ProtoReflect.Descriptor instead.
func (*CreateShopRequest) Descriptor() ([]byte, []int) {
	return file_shop_v1_shop_message_proto_rawDescGZIP(), []int{1}
}

func (x *CreateShopRequest) GetShop() *Shop {
	if x != nil {
		return x.Shop
	}
	return nil
}

type CreateShopResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Shop *Shop `protobuf:"bytes,1,opt,name=shop,proto3" json:"shop,omitempty"`
}

func (x *CreateShopResponse) Reset() {
	*x = CreateShopResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_shop_v1_shop_message_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateShopResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateShopResponse) ProtoMessage() {}

func (x *CreateShopResponse) ProtoReflect() protoreflect.Message {
	mi := &file_shop_v1_shop_message_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateShopResponse.ProtoReflect.Descriptor instead.
func (*CreateShopResponse) Descriptor() ([]byte, []int) {
	return file_shop_v1_shop_message_proto_rawDescGZIP(), []int{2}
}

func (x *CreateShopResponse) GetShop() *Shop {
	if x != nil {
		return x.Shop
	}
	return nil
}

var File_shop_v1_shop_message_proto protoreflect.FileDescriptor

var file_shop_v1_shop_message_proto_rawDesc = []byte{
	0x0a, 0x1a, 0x73, 0x68, 0x6f, 0x70, 0x2f, 0x76, 0x31, 0x2f, 0x73, 0x68, 0x6f, 0x70, 0x5f, 0x6d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x07, 0x73, 0x68,
	0x6f, 0x70, 0x2e, 0x76, 0x31, 0x22, 0x1f, 0x0a, 0x04, 0x53, 0x68, 0x6f, 0x70, 0x12, 0x17, 0x0a,
	0x07, 0x73, 0x68, 0x6f, 0x70, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06,
	0x73, 0x68, 0x6f, 0x70, 0x49, 0x64, 0x22, 0x36, 0x0a, 0x11, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65,
	0x53, 0x68, 0x6f, 0x70, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x21, 0x0a, 0x04, 0x73,
	0x68, 0x6f, 0x70, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0d, 0x2e, 0x73, 0x68, 0x6f, 0x70,
	0x2e, 0x76, 0x31, 0x2e, 0x53, 0x68, 0x6f, 0x70, 0x52, 0x04, 0x73, 0x68, 0x6f, 0x70, 0x22, 0x37,
	0x0a, 0x12, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x53, 0x68, 0x6f, 0x70, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x21, 0x0a, 0x04, 0x73, 0x68, 0x6f, 0x70, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x0d, 0x2e, 0x73, 0x68, 0x6f, 0x70, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x68, 0x6f,
	0x70, 0x52, 0x04, 0x73, 0x68, 0x6f, 0x70, 0x42, 0x0a, 0x5a, 0x08, 0x2e, 0x2f, 0x73, 0x68, 0x6f,
	0x70, 0x76, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_shop_v1_shop_message_proto_rawDescOnce sync.Once
	file_shop_v1_shop_message_proto_rawDescData = file_shop_v1_shop_message_proto_rawDesc
)

func file_shop_v1_shop_message_proto_rawDescGZIP() []byte {
	file_shop_v1_shop_message_proto_rawDescOnce.Do(func() {
		file_shop_v1_shop_message_proto_rawDescData = protoimpl.X.CompressGZIP(file_shop_v1_shop_message_proto_rawDescData)
	})
	return file_shop_v1_shop_message_proto_rawDescData
}

var file_shop_v1_shop_message_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_shop_v1_shop_message_proto_goTypes = []interface{}{
	(*Shop)(nil),               // 0: shop.v1.Shop
	(*CreateShopRequest)(nil),  // 1: shop.v1.CreateShopRequest
	(*CreateShopResponse)(nil), // 2: shop.v1.CreateShopResponse
}
var file_shop_v1_shop_message_proto_depIdxs = []int32{
	0, // 0: shop.v1.CreateShopRequest.shop:type_name -> shop.v1.Shop
	0, // 1: shop.v1.CreateShopResponse.shop:type_name -> shop.v1.Shop
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_shop_v1_shop_message_proto_init() }
func file_shop_v1_shop_message_proto_init() {
	if File_shop_v1_shop_message_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_shop_v1_shop_message_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Shop); i {
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
		file_shop_v1_shop_message_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateShopRequest); i {
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
		file_shop_v1_shop_message_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateShopResponse); i {
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
			RawDescriptor: file_shop_v1_shop_message_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_shop_v1_shop_message_proto_goTypes,
		DependencyIndexes: file_shop_v1_shop_message_proto_depIdxs,
		MessageInfos:      file_shop_v1_shop_message_proto_msgTypes,
	}.Build()
	File_shop_v1_shop_message_proto = out.File
	file_shop_v1_shop_message_proto_rawDesc = nil
	file_shop_v1_shop_message_proto_goTypes = nil
	file_shop_v1_shop_message_proto_depIdxs = nil
}
