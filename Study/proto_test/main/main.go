package main

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"littleShopee/Study/proto_test"
)

func main() {
	req := proto_test.HelloRequest{
		Name: "runzheng",
	}
	//序列化 引入 "github.com/golang/protobuf/proto"
	rsp, _ := proto.Marshal(&req)
	fmt.Println(string(rsp), rsp)

	//反序列化 先声明message这个结构体
	helloRequest := proto_test.HelloRequest{}
	_ = proto.Unmarshal(rsp, &helloRequest)
	fmt.Println(helloRequest.Name) //可以正常解码
}
