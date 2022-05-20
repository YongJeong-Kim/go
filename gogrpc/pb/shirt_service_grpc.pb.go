// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.19.3
// source: shirt_service.proto

package pb

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

// ShirtServiceClient is the client API for ShirtService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ShirtServiceClient interface {
	Broadcast(ctx context.Context, opts ...grpc.CallOption) (ShirtService_BroadcastClient, error)
}

type shirtServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewShirtServiceClient(cc grpc.ClientConnInterface) ShirtServiceClient {
	return &shirtServiceClient{cc}
}

func (c *shirtServiceClient) Broadcast(ctx context.Context, opts ...grpc.CallOption) (ShirtService_BroadcastClient, error) {
	stream, err := c.cc.NewStream(ctx, &ShirtService_ServiceDesc.Streams[0], "/ShirtService/Broadcast", opts...)
	if err != nil {
		return nil, err
	}
	x := &shirtServiceBroadcastClient{stream}
	return x, nil
}

type ShirtService_BroadcastClient interface {
	Send(*ShirtRequest) error
	Recv() (*ShirtResponse, error)
	grpc.ClientStream
}

type shirtServiceBroadcastClient struct {
	grpc.ClientStream
}

func (x *shirtServiceBroadcastClient) Send(m *ShirtRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *shirtServiceBroadcastClient) Recv() (*ShirtResponse, error) {
	m := new(ShirtResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// ShirtServiceServer is the server API for ShirtService service.
// All implementations must embed UnimplementedShirtServiceServer
// for forward compatibility
type ShirtServiceServer interface {
	Broadcast(ShirtService_BroadcastServer) error
	mustEmbedUnimplementedShirtServiceServer()
}

// UnimplementedShirtServiceServer must be embedded to have forward compatible implementations.
type UnimplementedShirtServiceServer struct {
}

func (UnimplementedShirtServiceServer) Broadcast(ShirtService_BroadcastServer) error {
	return status.Errorf(codes.Unimplemented, "method Broadcast not implemented")
}
func (UnimplementedShirtServiceServer) mustEmbedUnimplementedShirtServiceServer() {}

// UnsafeShirtServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ShirtServiceServer will
// result in compilation errors.
type UnsafeShirtServiceServer interface {
	mustEmbedUnimplementedShirtServiceServer()
}

func RegisterShirtServiceServer(s grpc.ServiceRegistrar, srv ShirtServiceServer) {
	s.RegisterService(&ShirtService_ServiceDesc, srv)
}

func _ShirtService_Broadcast_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(ShirtServiceServer).Broadcast(&shirtServiceBroadcastServer{stream})
}

type ShirtService_BroadcastServer interface {
	Send(*ShirtResponse) error
	Recv() (*ShirtRequest, error)
	grpc.ServerStream
}

type shirtServiceBroadcastServer struct {
	grpc.ServerStream
}

func (x *shirtServiceBroadcastServer) Send(m *ShirtResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *shirtServiceBroadcastServer) Recv() (*ShirtRequest, error) {
	m := new(ShirtRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// ShirtService_ServiceDesc is the grpc.ServiceDesc for ShirtService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ShirtService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "ShirtService",
	HandlerType: (*ShirtServiceServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Broadcast",
			Handler:       _ShirtService_Broadcast_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "shirt_service.proto",
}