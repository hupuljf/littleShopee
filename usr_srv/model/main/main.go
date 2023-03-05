package main

import (
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"littleShopee/usr_srv/handler"
	"littleShopee/usr_srv/proto"
	"net"
)

func main() {
	IP := flag.String("ip:", "0.0.0.0", "可以输入ip地址")
	Port := flag.Int("Port:", 50051, "可以输入端口号")
	flag.Parse()
	fmt.Println(*IP)
	server := grpc.NewServer()
	proto.RegisterUserServer(server, &handler.UserServer{})
	lis, _ := net.Listen("tcp", fmt.Sprintf("%s:%d", *IP, *Port))
	server.Serve(lis)

}
