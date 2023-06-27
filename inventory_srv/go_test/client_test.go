package go_test_test

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"littleShopee/inventory_srv/proto"
	"testing"
)

var client proto.InventoryClient
var conn *grpc.ClientConn

func TestInv(T *testing.T) {
	Init()
	//var x int32
	for i := 0; i <= 10; i++ {
		_, err := client.SetInv(context.Background(), &proto.GoodsInvInfo{
			GoodsId: int32(i),
			Num:     1000,
		})
		if err != nil {
			fmt.Println(err.Error())
		}
	}

	fmt.Println("设置库存成功")
	conn.Close()

}

func TestInvDetail(T *testing.T) {
	Init()
	rsp, err := client.InvDetail(context.Background(), &proto.GoodsInvInfo{
		GoodsId: 421,
	})
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(rsp)
	conn.Close()

}

func Init() {
	conn, err := grpc.Dial("192.168.2.3:63874", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	client = proto.NewInventoryClient(conn)

}
