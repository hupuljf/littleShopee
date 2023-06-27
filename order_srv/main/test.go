package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"littleShopee/order_srv/proto"
)

var orderClient proto.OrderClient
var conn *grpc.ClientConn

func Init() {
	conn, err := grpc.Dial("192.168.2.3:54880", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	orderClient = proto.NewOrderClient(conn)

}

func CreateCartItem(userId, nums, goodsId int32) {
	rsp, err := orderClient.CreateCartItem(context.Background(), &proto.CartItemRequest{
		UserId:  userId,
		Nums:    nums,
		GoodsId: goodsId,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(rsp.Id)
}

func CartItemList(userId int32) {
	rsp, err := orderClient.CartItemList(context.Background(), &proto.UserInfo{
		Id: userId,
	})
	if err != nil {
		panic(err)
	}
	for _, item := range rsp.Data {
		fmt.Println(item.Id, item.GoodsId, item.Nums)
	}
}

func UpdateCartItem(uid int32, gid int32) {
	_, err := orderClient.UpdateCartItem(context.Background(), &proto.CartItemRequest{
		UserId:  uid,
		GoodsId: gid,
		Checked: true,
	})
	if err != nil {
		panic(err)
	}
}

func CreateOrder() {
	_, err := orderClient.CreateOrder(context.Background(), &proto.OrderRequest{
		UserId:  1,
		Address: "北京市",
		Name:    "bobby",
		Mobile:  "18787878787",
		Post:    "请尽快发货",
	})
	if err != nil {
		panic(err)
	}
}

func main() {
	Init()
	//CreateCartItem(1, 1, 422)
	//CartItemList(1)
	//UpdateCartItem(1, 422)

	CreateOrder()
	err := conn.Close()
	fmt.Println(err)

}
