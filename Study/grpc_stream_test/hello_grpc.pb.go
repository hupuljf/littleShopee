// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.12
// source: hello.proto

package grpc_stream_test

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

// HelloClient is the client API for Hello service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type HelloClient interface {
	HelloServerStream(ctx context.Context, in *HelloRequest, opts ...grpc.CallOption) (Hello_HelloServerStreamClient, error)
	HelloClientStream(ctx context.Context, opts ...grpc.CallOption) (Hello_HelloClientStreamClient, error)
	HelloEachStream(ctx context.Context, opts ...grpc.CallOption) (Hello_HelloEachStreamClient, error)
}

type helloClient struct {
	cc grpc.ClientConnInterface
}

func NewHelloClient(cc grpc.ClientConnInterface) HelloClient {
	return &helloClient{cc}
}

func (c *helloClient) HelloServerStream(ctx context.Context, in *HelloRequest, opts ...grpc.CallOption) (Hello_HelloServerStreamClient, error) {
	stream, err := c.cc.NewStream(ctx, &Hello_ServiceDesc.Streams[0], "/Hello/HelloServerStream", opts...)
	if err != nil {
		return nil, err
	}
	x := &helloHelloServerStreamClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Hello_HelloServerStreamClient interface {
	Recv() (*HelloResponse, error)
	grpc.ClientStream
}

type helloHelloServerStreamClient struct {
	grpc.ClientStream
}

func (x *helloHelloServerStreamClient) Recv() (*HelloResponse, error) {
	m := new(HelloResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *helloClient) HelloClientStream(ctx context.Context, opts ...grpc.CallOption) (Hello_HelloClientStreamClient, error) {
	stream, err := c.cc.NewStream(ctx, &Hello_ServiceDesc.Streams[1], "/Hello/HelloClientStream", opts...)
	if err != nil {
		return nil, err
	}
	x := &helloHelloClientStreamClient{stream}
	return x, nil
}

type Hello_HelloClientStreamClient interface {
	Send(*HelloRequest) error
	CloseAndRecv() (*HelloResponse, error)
	grpc.ClientStream
}

type helloHelloClientStreamClient struct {
	grpc.ClientStream
}

func (x *helloHelloClientStreamClient) Send(m *HelloRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *helloHelloClientStreamClient) CloseAndRecv() (*HelloResponse, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(HelloResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *helloClient) HelloEachStream(ctx context.Context, opts ...grpc.CallOption) (Hello_HelloEachStreamClient, error) {
	stream, err := c.cc.NewStream(ctx, &Hello_ServiceDesc.Streams[2], "/Hello/HelloEachStream", opts...)
	if err != nil {
		return nil, err
	}
	x := &helloHelloEachStreamClient{stream}
	return x, nil
}

type Hello_HelloEachStreamClient interface {
	Send(*HelloRequest) error
	Recv() (*HelloResponse, error)
	grpc.ClientStream
}

type helloHelloEachStreamClient struct {
	grpc.ClientStream
}

func (x *helloHelloEachStreamClient) Send(m *HelloRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *helloHelloEachStreamClient) Recv() (*HelloResponse, error) {
	m := new(HelloResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// HelloServer is the server API for Hello service.
// All implementations must embed UnimplementedHelloServer
// for forward compatibility
type HelloServer interface {
	HelloServerStream(*HelloRequest, Hello_HelloServerStreamServer) error
	HelloClientStream(Hello_HelloClientStreamServer) error
	HelloEachStream(Hello_HelloEachStreamServer) error
	//mustEmbedUnimplementedHelloServer()
}

// UnimplementedHelloServer must be embedded to have forward compatible implementations.
type UnimplementedHelloServer struct {
}

func (UnimplementedHelloServer) HelloServerStream(*HelloRequest, Hello_HelloServerStreamServer) error {
	return status.Errorf(codes.Unimplemented, "method HelloServerStream not implemented")
}
func (UnimplementedHelloServer) HelloClientStream(Hello_HelloClientStreamServer) error {
	return status.Errorf(codes.Unimplemented, "method HelloClientStream not implemented")
}
func (UnimplementedHelloServer) HelloEachStream(Hello_HelloEachStreamServer) error {
	return status.Errorf(codes.Unimplemented, "method HelloEachStream not implemented")
}
func (UnimplementedHelloServer) mustEmbedUnimplementedHelloServer() {}

// UnsafeHelloServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to HelloServer will
// result in compilation errors.
type UnsafeHelloServer interface {
	mustEmbedUnimplementedHelloServer()
}

func RegisterHelloServer(s grpc.ServiceRegistrar, srv HelloServer) {
	s.RegisterService(&Hello_ServiceDesc, srv)
}

func _Hello_HelloServerStream_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(HelloRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(HelloServer).HelloServerStream(m, &helloHelloServerStreamServer{stream})
}

type Hello_HelloServerStreamServer interface {
	Send(*HelloResponse) error
	grpc.ServerStream
}

type helloHelloServerStreamServer struct {
	grpc.ServerStream
}

func (x *helloHelloServerStreamServer) Send(m *HelloResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _Hello_HelloClientStream_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(HelloServer).HelloClientStream(&helloHelloClientStreamServer{stream})
}

type Hello_HelloClientStreamServer interface {
	SendAndClose(*HelloResponse) error
	Recv() (*HelloRequest, error)
	grpc.ServerStream
}

type helloHelloClientStreamServer struct {
	grpc.ServerStream
}

func (x *helloHelloClientStreamServer) SendAndClose(m *HelloResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *helloHelloClientStreamServer) Recv() (*HelloRequest, error) {
	m := new(HelloRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _Hello_HelloEachStream_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(HelloServer).HelloEachStream(&helloHelloEachStreamServer{stream})
}

type Hello_HelloEachStreamServer interface {
	Send(*HelloResponse) error
	Recv() (*HelloRequest, error)
	grpc.ServerStream
}

type helloHelloEachStreamServer struct {
	grpc.ServerStream
}

func (x *helloHelloEachStreamServer) Send(m *HelloResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *helloHelloEachStreamServer) Recv() (*HelloRequest, error) {
	m := new(HelloRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// Hello_ServiceDesc is the grpc.ServiceDesc for Hello service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Hello_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "Hello",
	HandlerType: (*HelloServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "HelloServerStream",
			Handler:       _Hello_HelloServerStream_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "HelloClientStream",
			Handler:       _Hello_HelloClientStream_Handler,
			ClientStreams: true,
		},
		{
			StreamName:    "HelloEachStream",
			Handler:       _Hello_HelloEachStream_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "hello.proto",
}
