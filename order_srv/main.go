package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/hashicorp/consul/api"
	"github.com/satori/go.uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"littleShopee/order_srv/global"
	"littleShopee/order_srv/handler"
	"littleShopee/order_srv/initialize"
	"littleShopee/order_srv/proto"
	"littleShopee/order_srv/utils"
	"net"
)

func main() {

	IP := flag.String("ip:", "0.0.0.0", "可以输入ip地址")
	Port := flag.Int("Port:", 0, "可以输入端口号")
	flag.Parse()
	zap.S().Info(*IP)
	initialize.InitLogger()
	initialize.InitConfig()
	initialize.InitDB()
	initialize.InitSrvs()
	server := grpc.NewServer()
	proto.RegisterOrderServer(server, &handler.OrderServer{})
	//非输入情况下 修改端口 以随机端口代替
	if *Port == 0 {
		*Port, _ = utils.GetFreePort()
	}
	lis, _ := net.Listen("tcp", fmt.Sprintf("%s:%d", *IP, *Port))

	//grpc的健康检查
	grpc_health_v1.RegisterHealthServer(server, health.NewServer())
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", global.ServerConfig.ConsulInfo.Host, global.ServerConfig.ConsulInfo.Port) //consul的出发ip和端口
	//在consul上注册服务
	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}
	//生成对应的检查对象
	check := &api.AgentServiceCheck{
		GRPC:                           fmt.Sprintf("192.168.2.3:%d", *Port), //得用本机地址
		Timeout:                        "5s",
		Interval:                       "5s",
		DeregisterCriticalServiceAfter: "60s",
	}
	//生成注册对象
	serviceID := fmt.Sprintf("%s", uuid.NewV4())
	registration := new(api.AgentServiceRegistration)
	registration.Name = global.ServerConfig.Name
	registration.ID = serviceID
	registration.Port = *Port
	registration.Tags = []string{"grpc", "order"}
	registration.Address = "192.168.2.3"
	registration.Check = check

	err = client.Agent().ServiceRegister(registration)
	//client.Agent().ServiceDeregister() //注销服务
	if err != nil {
		panic(err)
	}
	go func() {
		server.Serve(lis)
	}() //不在协程里的话 后面的执行不到

	//接受终止信号 可以优雅退出
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit //quit这个通道拿到东西前都阻塞在这里
	if err = client.Agent().ServiceDeregister(serviceID); err != nil {
		zap.S().Info("注销失败")
		panic(err)
	}

}