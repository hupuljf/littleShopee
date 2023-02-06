package main

import (
	"fmt"
	"net/rpc"
)

func main() {
	//建立连接 打到对方的IP和端口
	client, e := rpc.Dial("tcp", "localhost:8000")
	fmt.Println("1", e)
	var reply string
	err := client.Call("HelloService.Hello", "runzheng", &reply)
	fmt.Println(err)
	fmt.Println(reply)
	//var reply *string = new(string) //得先给这个指针赋予空间 不然传进去的就是空指针
	//_ = client.Call("HelloService.Hello", "runzheng", reply)
	//fmt.Println(*reply)

}
