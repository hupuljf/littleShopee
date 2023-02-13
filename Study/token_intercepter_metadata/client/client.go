package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"littleShopee/Study/proto_test"
	"time"
)

func main() {
	//conn, _ := grpc.Dial("0.0.0.0:9091", grpc.WithInsecure())
	////程序退出需要关闭链接
	//defer conn.Close()
	//
	////new 一个grpc的client
	//c := proto_test.NewHelloClient(conn)
	////这个hello就是相当于server的处理函数 需要将需要参数传入（后台再传到server端）
	////metadata
	////md := metadata.Pairs("token", "1111jijio")
	//md := metadata.New(map[string]string{
	//	"token": "ggg000",
	//	"name":  "bbb",
	//})
	//ctx := metadata.NewOutgoingContext(context.Background(), md)
	//r, _ := c.Hello(ctx, &proto_test.HelloRequest{
	//	Name:    "jiba",
	//	Age:     0,
	//	Courses: nil,
	//})
	//fmt.Println(r.Reply)

	/////// 带拦截器的 client
	interceptor := func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		start := time.Now()
		//客户端建立连接后发起请求之前做的一些事情写在这儿
		//这个opts是调用前的选择
		//引入metadata 相当于客户端的请求带header
		md := metadata.New(map[string]string{
			"appid":  "1001",
			"appkey": "1001key",
		})
		//将metadata带入ctx
		ctx = metadata.NewOutgoingContext(ctx, md)
		err := invoker(ctx, method, req, reply, cc, opts...)
		fmt.Printf("耗时：%s\n", time.Since(start))
		return err
	}
	//这个option是建立连接前的选择
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	opts = append(opts, grpc.WithUnaryInterceptor(interceptor))
	conn, err := grpc.Dial("127.0.0.1:9091", opts...)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	c := proto_test.NewHelloClient(conn)
	r, err := c.Hello(context.Background(), &proto_test.HelloRequest{Name: "bobby"})
	if err != nil {
		panic(err)
	}
	fmt.Println(r.Reply)

}
