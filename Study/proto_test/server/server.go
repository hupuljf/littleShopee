package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"net"

	"littleShopee/Study/proto_test"
)

//实现server接口
type Server struct{}

func (s *Server) Hello(ctx context.Context, request *proto_test.HelloRequest) (*proto_test.Response, error) {
	//ctx.Err()
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		//打印的是一个string类型
		fmt.Println(md["token"])
	}

	return &proto_test.Response{
		Reply: "Mother Fucker! " + (*request).Name,
	}, nil

}

//手动将这个函数改为大写 才可以调用 实际上不知道这个函数做什么的
func (s *Server) MustEmbedUnimplementedHelloServer() {

}

func main() {
	//正常server端写法
	//g := grpc.NewServer()
	//proto_test.RegisterHelloServer(g, &Server{})

	//加了拦截器的server端写法
	//1.构造一个拦截器 是一个函数
	interceptor := func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		fmt.Println("接收到了一个新的请求，处理一些在这个请求之前的事情")
		//处理请求直接在拦截器里面做了
		res, err := handler(ctx, req)
		fmt.Println("请求已经完成")
		return res, err
	}
	//opt有多个 可以有多个拦截器
	opt := grpc.UnaryInterceptor(interceptor)
	g := grpc.NewServer(opt)
	proto_test.RegisterHelloServer(g, &Server{})

	//grpc服务的监听端口设置
	listener, _ := net.Listen("tcp", "0.0.0.0:9091")
	g.Serve(listener)
}
