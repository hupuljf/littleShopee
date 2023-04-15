package main

import (
	"context"
	"encoding/json"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"littleShopee/goods_srv/model"
	"littleShopee/goods_srv/proto"
)

var userClient proto.GoodsClient
var conn *grpc.ClientConn

func Init() {
	var err error
	conn, err = grpc.Dial("127.0.0.1:50051", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	userClient = proto.NewGoodsClient(conn)
}

func TestGetUserList() {
	rsp, err := userClient.BrandList(context.Background(), &proto.BrandFilterRequest{
		Pages:       1,
		PagePerNums: 5,
	})
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	fmt.Println(rsp)
}

func TestCategory() {
	rsp, err := userClient.GetAllCategorysList(context.Background(), &emptypb.Empty{})
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	b := []byte(rsp.JsonData)
	var categorys []model.Category
	json.Unmarshal(b, &categorys)
	fmt.Println(rsp)
	fmt.Println(categorys[0].SubCategory[0])
}

func TestCategoryBrandList() {
	rsp, err := userClient.CategoryBrandList(context.Background(), &proto.CategoryBrandFilterRequest{
		PagePerNums: 10,
		Pages:       1,
	})
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	fmt.Println(rsp.Data[0])

}
func TestBrandList() {
	rsp, err := userClient.BrandList(context.Background(), &proto.BrandFilterRequest{
		PagePerNums: 10,
		Pages:       1,
	})
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	fmt.Println(rsp)

}

//func TestCreateUser() {
//	for i := 0; i < 10; i++ {
//		rsp, err := userClient.CreateUser(context.Background(), &proto.CreateUserInfo{
//			NickName: fmt.Sprintf("bobby%d", i),
//			Mobile:   fmt.Sprintf("1371949551%d", i),
//			PassWord: "admin123" + strconv.Itoa(i),
//		})
//		if err != nil {
//			panic(err)
//		}
//		fmt.Println(rsp.Id)
//	}
//}

func TestGoodsList() {
	rsp, err := userClient.GoodsList(context.Background(), &proto.GoodsFilterRequest{
		PagePerNums: 10,
		Pages:       1,
	})
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	fmt.Println(rsp.Data[0])

}

func main() {
	//userList := new(proto.UserListResponse)
	//(*userList).Total = int32(999)
	//fmt.Println(userList.Total)
	var err error
	//建立连接 dial ip端口
	conn, err = grpc.Dial("0.0.0.0:50051", grpc.WithInsecure())
	defer conn.Close()
	if err != nil {
		fmt.Println("jjjjjjjjjjjjj")
		panic(err)
	}
	userClient = proto.NewGoodsClient(conn)
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
	TestGoodsList()
}
