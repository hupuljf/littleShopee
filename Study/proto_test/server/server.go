package main

import (
	"context"
	"net"

	"google.golang.org/grpc"
	"littleShopee/Study/proto_test"
)

//实现server接口
type Server struct{}

func (s *Server) Hello(ctx context.Context, request *proto_test.HelloRequest) (*proto_test.Response, error) {
	//ctx.Err()
	return &proto_test.Response{
		Reply: "Mother Fucker! " + (*request).Name,
	}, nil

}

//手动将这个函数改为大写 才可以调用 实际上不知道这个函数做什么的
func (s *Server) MustEmbedUnimplementedHelloServer() {

}

func main() {
	g := grpc.NewServer()
	proto_test.RegisterHelloServer(g, &Server{})
	//grpc服务的监听端口设置
	listener, _ := net.Listen("tcp", "0.0.0.0:9090")
	g.Serve(listener)
}
