package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"littleShopee/Study/proto_test"
)

func main() {
	conn, _ := grpc.Dial("0.0.0.0:9091", grpc.WithInsecure())
	//程序退出需要关闭链接
	defer conn.Close()

	//new 一个grpc的client
	c := proto_test.NewHelloClient(conn)
	//这个hello就是相当于server的处理函数 需要将需要参数传入（后台再传到server端）
	//metadata
	//md := metadata.Pairs("token", "1111jijio")
	md := metadata.New(map[string]string{
		"token": "ggg000",
		"name":  "bbb",
	})
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	r, _ := c.Hello(ctx, &proto_test.HelloRequest{
		Name:    "jiba",
		Age:     0,
		Courses: nil,
	})
	fmt.Println(r.Reply)

}
