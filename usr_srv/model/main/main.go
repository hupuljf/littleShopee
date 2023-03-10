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
	//传入一批user
	//users := make([]model.User, 0)
	//for i := 0; i < 10; i++ {
	//	name := "test-" + strconv.Itoa(i)
	//	users = append(users, model.User{
	//		NickName: name,
	//		Mobile:   "1371959451" + strconv.Itoa(i),
	//		Password: "admin123" + strconv.Itoa(i),
	//	})
	//}
	//global.DB.Create(&users)

	IP := flag.String("ip:", "0.0.0.0", "可以输入ip地址")
	Port := flag.Int("Port:", 50051, "可以输入端口号")
	flag.Parse()
	fmt.Println(*IP)
	server := grpc.NewServer()
	proto.RegisterUserServer(server, &handler.UserServer{})
	lis, _ := net.Listen("tcp", fmt.Sprintf("%s:%d", *IP, *Port))
	server.Serve(lis)

}
