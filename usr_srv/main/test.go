package main

import (
	"context"
	"fmt"
	"littleShopee/usr_srv/proto"
	"strconv"

	"google.golang.org/grpc"
)

var userClient proto.UserClient
var conn *grpc.ClientConn

func Init() {
	var err error
	conn, err = grpc.Dial("127.0.0.1:50051", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	userClient = proto.NewUserClient(conn)
}

func TestGetUserList() {
	rsp, err := userClient.GetUserList(context.Background(), &proto.PageInfo{
		Pn:    1,
		PSize: 5,
	})
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	for _, user := range rsp.Data {
		fmt.Println(user.Mobile, user.NickName, user.PassWord)
		checkRsp, err := userClient.CheckPassWord(context.Background(), &proto.PasswordCheckInfo{
			Password:          "admin123",
			EncryptedPassword: user.PassWord,
		})
		if err != nil {
			panic(err)
		}
		fmt.Println(checkRsp.Success)
	}
}

func TestCreateUser() {
	for i := 0; i < 10; i++ {
		rsp, err := userClient.CreateUser(context.Background(), &proto.CreateUserInfo{
			NickName: fmt.Sprintf("bobby%d", i),
			Mobile:   fmt.Sprintf("1371949551%d", i),
			PassWord: "admin123" + strconv.Itoa(i),
		})
		if err != nil {
			panic(err)
		}
		fmt.Println(rsp.Id)
	}
}

func main() {
	userList := new(proto.UserListResponse)
	(*userList).Total = int32(999)
	fmt.Println(userList.Total)
	var err error
	//建立连接 dial ip端口
	conn, err = grpc.Dial("0.0.0.0:50051", grpc.WithInsecure())
	defer conn.Close()
	if err != nil {
		fmt.Println("jjjjjjjjjjjjj")
		panic(err)
	}
	userClient = proto.NewUserClient(conn)
	//for i := 0; i < 10; i++ {
	//	rsp, err := userClient.CreateUser(context.Background(), &proto.CreateUserInfo{
	//		NickName: fmt.Sprintf("bobby%d", i),
	//		Mobile:   fmt.Sprintf("1371949551%d", i),
	//		PassWord: "admin123" + strconv.Itoa(i),
	//	})
	//	if err != nil {
	//		panic(err)
	//	}
	//	fmt.Println(rsp.Id)
	//}

	TestGetUserList()

}
