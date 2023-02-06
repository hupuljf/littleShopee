package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"littleShopee/Study/proto_test"
)

func main() {
	conn, _ := grpc.Dial("0.0.0.0:9090", grpc.WithInsecure())
	//程序退出需要关闭链接
	defer conn.Close()

	//new 一个grpc的client
	c := proto_test.NewHelloClient(conn)
	//这个hello就是相当于server的处理函数 需要将需要参数传入（后台再传到server端）
	r, _ := c.Hello(context.Background(), &proto_test.HelloRequest{
		Name:    "jiba",
		Age:     0,
		Courses: nil,
	})
	fmt.Println(r.Reply)

}
