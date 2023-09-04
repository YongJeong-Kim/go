// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             (unknown)
// source: user/v1/user_service.proto

package userv1

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

// SimpleServerClient is the client API for SimpleServer service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SimpleServerClient interface {
	CreateUser(ctx context.Context, in *CreateUserRequest, opts ...grpc.CallOption) (*CreateUserResponse, error)
	DeleteUser(ctx context.Context, in *DeleteUserRequest, opts ...grpc.CallOption) (*DeleteUserResponse, error)
	//  rpc UploadUser(UploadUserRequest) returns (google.protobuf.Empty) {
	UploadUser(ctx context.Context, in *UploadUserRequest, opts ...grpc.CallOption) (*UploadUserResponse, error)
}

type simpleServerClient struct {
	cc grpc.ClientConnInterface
}

func NewSimpleServerClient(cc grpc.ClientConnInterface) SimpleServerClient {
	return &simpleServerClient{cc}
}

func (c *simpleServerClient) CreateUser(ctx context.Context, in *CreateUserRequest, opts ...grpc.CallOption) (*CreateUserResponse, error) {
	out := new(CreateUserResponse)
	err := c.cc.Invoke(ctx, "/user.v1.SimpleServer/CreateUser", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *simpleServerClient) DeleteUser(ctx context.Context, in *DeleteUserRequest, opts ...grpc.CallOption) (*DeleteUserResponse, error) {
	out := new(DeleteUserResponse)
	err := c.cc.Invoke(ctx, "/user.v1.SimpleServer/DeleteUser", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *simpleServerClient) UploadUser(ctx context.Context, in *UploadUserRequest, opts ...grpc.CallOption) (*UploadUserResponse, error) {
	out := new(UploadUserResponse)
	err := c.cc.Invoke(ctx, "/user.v1.SimpleServer/UploadUser", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SimpleServerServer is the server API for SimpleServer service.
// All implementations must embed UnimplementedSimpleServerServer
// for forward compatibility
type SimpleServerServer interface {
	CreateUser(context.Context, *CreateUserRequest) (*CreateUserResponse, error)
	DeleteUser(context.Context, *DeleteUserRequest) (*DeleteUserResponse, error)
	//  rpc UploadUser(UploadUserRequest) returns (google.protobuf.Empty) {
	UploadUser(context.Context, *UploadUserRequest) (*UploadUserResponse, error)
	mustEmbedUnimplementedSimpleServerServer()
}

// UnimplementedSimpleServerServer must be embedded to have forward compatible implementations.
type UnimplementedSimpleServerServer struct {
}

func (UnimplementedSimpleServerServer) CreateUser(context.Context, *CreateUserRequest) (*CreateUserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateUser not implemented")
}
func (UnimplementedSimpleServerServer) DeleteUser(context.Context, *DeleteUserRequest) (*DeleteUserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteUser not implemented")
}
func (UnimplementedSimpleServerServer) UploadUser(context.Context, *UploadUserRequest) (*UploadUserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UploadUser not implemented")
}
func (UnimplementedSimpleServerServer) mustEmbedUnimplementedSimpleServerServer() {}

// UnsafeSimpleServerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SimpleServerServer will
// result in compilation errors.
type UnsafeSimpleServerServer interface {
	mustEmbedUnimplementedSimpleServerServer()
}

func RegisterSimpleServerServer(s grpc.ServiceRegistrar, srv SimpleServerServer) {
	s.RegisterService(&SimpleServer_ServiceDesc, srv)
}

func _SimpleServer_CreateUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SimpleServerServer).CreateUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/user.v1.SimpleServer/CreateUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SimpleServerServer).CreateUser(ctx, req.(*CreateUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SimpleServer_DeleteUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SimpleServerServer).DeleteUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/user.v1.SimpleServer/DeleteUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SimpleServerServer).DeleteUser(ctx, req.(*DeleteUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SimpleServer_UploadUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UploadUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SimpleServerServer).UploadUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/user.v1.SimpleServer/UploadUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SimpleServerServer).UploadUser(ctx, req.(*UploadUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// SimpleServer_ServiceDesc is the grpc.ServiceDesc for SimpleServer service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var SimpleServer_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "user.v1.SimpleServer",
	HandlerType: (*SimpleServerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateUser",
			Handler:    _SimpleServer_CreateUser_Handler,
		},
		{
			MethodName: "DeleteUser",
			Handler:    _SimpleServer_DeleteUser_Handler,
		},
		{
			MethodName: "UploadUser",
			Handler:    _SimpleServer_UploadUser_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "user/v1/user_service.proto",
}
