package main

import (
	"net"
	"net/rpc"
)

//定义一个有函数的接口（rpc服务注册函数的传参是接口）
type HelloService struct {
}

func (s *HelloService) Hello(request string, reply *string) error {
	*reply = "Hello" + request
	return nil
}

func main() {
	//1.实例化一个server 使用listner监听端口
	listener, _ := net.Listen("tcp", ":8000")
	//2.注册一个rpc服务
	_ = rpc.RegisterName("HelloService", &HelloService{})
	//3.建立连接 启动服务
	conn, _ := listener.Accept()
	rpc.ServeConn(conn) //把这个连接的处理给rpc
}
