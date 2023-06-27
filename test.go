package main

import (
	"context"
	"google.golang.org/grpc"
	"littleShopee/inventory_srv/proto"
	"sync"
)

func main() {
	//dsn := "root:8971841xm@tcp(localhost:3306)/shop_inventory_srv?charset=utf8mb4&parseTime=True&loc=Local"
	//newLogger := logger.New(
	//	log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
	//	logger.Config{
	//		SlowThreshold: time.Second, // 慢 SQL 阈值
	//		LogLevel:      logger.Info, // Log level
	//		Colorful:      true,        // 禁用彩色打印
	//	},
	//)
	//// 全局模式
	//db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
	//	Logger: newLogger,
	//}) //初始化gorm 给db赋值 并且初始化gorm该有的配置
	//if err != nil {
	//	panic(err)
	//}
	conn, err := grpc.Dial("192.168.2.3:53098", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	client := proto.NewInventoryClient(conn)

	//测试分布式情况下不加锁的结果
	wg := sync.WaitGroup{}
	wg.Add(20)
	for i := 0; i < 20; i++ {
		go func() {
			defer wg.Done()
			client.Sell(context.Background(), &proto.SellInfo{
				GoodsInfo: []*proto.GoodsInvInfo{{
					GoodsId: 445,
					Num:     10,
				},
				},
			})
		}()
	}
	wg.Wait()

}
