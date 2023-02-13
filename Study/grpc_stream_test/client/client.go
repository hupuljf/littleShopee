package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"littleShopee/Study/grpc_stream_test"
	"sync"
	"time"
)

var clientRpc grpc_stream_test.HelloClient

//服务端流模式的请求 client请求 server源源不断
func serverStreamDemo() {
	//服务端流模式
	res, err := clientRpc.HelloServerStream(context.Background(), &grpc_stream_test.HelloRequest{Name: "Fuck! server"}) //context和客户端的请求体
	if err != nil {
		panic("rpc请求错误：" + err.Error())
	}
	for {
		data, err := res.Recv() //是reply
		if err != nil {
			fmt.Println("客户端发送完了:", err)
			return
		}
		fmt.Println("客户端返回数据流值:", data.Reply)
	}
}

//客户端流模式请求
func clientStreamDemo() {
	//客户端流模式
	cliStr, err := clientRpc.HelloClientStream(context.Background())
	if err != nil {
		panic("rpc请求错误：" + err.Error())
	}
	i := 0
	for {
		i++
		_ = cliStr.Send(&grpc_stream_test.HelloRequest{
			Name: "Fuck",
		})
		time.Sleep(time.Second * 1)
		if i > 10 {
			break
		}
	}
}

//双向流模式请求
func clientAndServerStreamDemo() {
	//双向流模式
	allStr, _ := clientRpc.HelloEachStream(context.Background())
	wg := sync.WaitGroup{}
	wg.Add(1)
	//接受服务端消息的协程
	go func() {
		defer wg.Done()
		for {
			//业务代码
			res, err := allStr.Recv()
			if err != nil {
				fmt.Println("本次服务端流数据发送完了:", err)
				break
			}
			fmt.Println("收到服务端发来消息：", res.Reply)
		}
	}()

	//发送消息给服务端的协程 wg值为1 这个协程结束就都结束了
	go func() {
		defer wg.Done()
		i := 0
		for {
			i++
			//业务代码
			_ = allStr.Send(&grpc_stream_test.HelloRequest{
				Name: "Fuck Each",
			})
			time.Sleep(time.Second * 1)
			if i > 10 {
				break
			}
		}
	}()
	wg.Wait()
}

func main() {
	conn, err := grpc.Dial("127.0.0.1:9090", grpc.WithInsecure())
	if err != nil {
		panic("rpc连接错误：" + err.Error())
	}
	defer conn.Close()
	//new一个client绑定连接
	clientRpc = grpc_stream_test.NewHelloClient(conn)
	//serverStreamDemo()
	clientAndServerStreamDemo()
}
