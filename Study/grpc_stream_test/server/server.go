package main

import (
	"fmt"
	"littleShopee/Study/grpc_stream_test"
	"net"
	"sync"
	"time"

	"google.golang.org/grpc"
)

//type HelloServer interface {
//	HelloServerStream(*HelloRequest, Hello_HelloServerStreamServer) error
//	HelloClientStream(Hello_HelloClientStreamServer) error
//	HelloEachStream(Hello_HelloEachStreamServer) error
//	//mustEmbedUnimplementedHelloServer()
//}
//实现server接口
type Server struct{}

//grpc_stream_test.Hello_HelloServerStreamServer 这个接口用来发送响应 可持续发送
func (s *Server) HelloServerStream(req *grpc_stream_test.HelloRequest, res grpc_stream_test.Hello_HelloServerStreamServer) error {
	i := 0
	for {
		i++
		//业务代码 客户端一个请求 服务端源源不断地把响应send出来
		_ = res.Send(&grpc_stream_test.HelloResponse{
			Reply: fmt.Sprintf("这是发给%s的数据流", req.Name),
		})
		time.Sleep(time.Second * 1)
		if i > 10 {
			break
		}
	}
	return nil
}

////手动将这个函数改为大写 才可以调用 实际上不知道这个函数做什么的
//func (s *Server) MustEmbedUnimplementedHelloServer() {
//
//}

//客户端流模式
func (s *Server) HelloClientStream(cliStr grpc_stream_test.Hello_HelloClientStreamServer) error {
	for {
		//业务代码 服务端只是源源不断地接受客户端地数据
		res, err := cliStr.Recv()
		if err != nil {
			fmt.Println("本次客户端流数据发送完了:", err)
			break
		}
		fmt.Println("客户端发来消息：", res.Name) //res 返回的是helloRequest（客户端的数据）
	}
	return nil
}

//双向流模式 搞了两个协程 把上面两种情况都加上
func (s *Server) HelloEachStream(allStr grpc_stream_test.Hello_HelloEachStreamServer) error {
	wg := sync.WaitGroup{}
	wg.Add(2)
	//接受客户端消息的协程
	go func() {
		defer wg.Done()
		for {
			//业务代码
			res, err := allStr.Recv()
			if err != nil {
				fmt.Println("本次客户端流数据发送完了:", err)
				break
			}
			fmt.Println("收到客户端发来消息：", res.Name)
		}
	}()

	//发送消息给客户端的协程
	go func() {
		defer wg.Done()
		i := 0
		for {
			i++
			//业务代码
			_ = allStr.Send(&grpc_stream_test.HelloResponse{
				Reply: fmt.Sprintf("这是发给客户端的数据流"), //没有hellorequest了 数据流就是上面拿到的数据
			})
			time.Sleep(time.Second * 1)
			if i > 10 {
				break
			}
		}
	}()
	wg.Wait()
	return nil
}

func main() {
	g := grpc.NewServer()
	grpc_stream_test.RegisterHelloServer(g, &Server{})
	//grpc服务的监听端口设置
	listener, _ := net.Listen("tcp", "0.0.0.0:9090")
	g.Serve(listener)
}
