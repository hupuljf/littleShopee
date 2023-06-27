package handler

import (
	"context"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
	"littleShopee/inventory_srv/global"
	"littleShopee/inventory_srv/model"

	//"google.golang.org/protobuf/runtime/protoimpl"
	"littleShopee/inventory_srv/proto"
)

type InventoryServer struct {
	proto.UnimplementedInventoryServer
}

//type BrandFilterRequest struct {
//	state         protoimpl.MessageState
//	sizeCache     protoimpl.SizeCache
//	unknownFields protoimpl.UnknownFields
//
//	Pages       int32 `protobuf:"varint,1,opt,name=pages,proto3" json:"pages,omitempty"`
//	PagePerNums int32 `protobuf:"varint,2,opt,name=pagePerNums,proto3" json:"pagePerNums,omitempty"`
//}

//分页函数
func Paginate(pn, ps int) func(db *gorm.DB) *gorm.DB { //返回的是一个函数 这个函数返回的是gorm.db
	return func(db *gorm.DB) *gorm.DB {
		if pn == 0 {
			pn = 1
		}

		offset := (pn - 1) * ps
		return db.Offset(offset).Limit(ps)
	}
}

func (InventoryServer) SetInv(ctx context.Context, req *proto.GoodsInvInfo) (*emptypb.Empty, error) {
	//设置库存 更新库存 需要先查有没有库存
	var inv model.Inventory
	global.DB.Where(&model.Inventory{Goods: req.GoodsId}).First(&inv) //查到了就更新 查不到 inv是nil 直接新建
	inv.Goods = req.GoodsId
	inv.Stocks = req.Num

	global.DB.Save(&inv)
	return &emptypb.Empty{}, nil

}

func (InventoryServer) InvDetail(ctx context.Context, req *proto.GoodsInvInfo) (*proto.GoodsInvInfo, error) {
	var inv model.Inventory
	if res := global.DB.Where(&model.Inventory{
		Goods: req.GoodsId,
	}).First(&inv); res.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "该商品暂无库存信息")
	}
	return &proto.GoodsInvInfo{
		GoodsId: inv.Goods,
		Num:     inv.Stocks,
	}, nil
}

//func (InventoryServer) Sell(ctx context.Context, req *proto.SellInfo) (*emptypb.Empty, error) {
//	//得开启一个事务 一批商品要么全成功 要么全失败
//	//后续需要加分布书锁
//	tx := global.DB.Begin()
//	for _, goodInfo := range req.GoodsInfo {
//		var inv model.Inventory
//		if res := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where(&model.Inventory{ //悲观锁
//			Goods: goodInfo.GoodsId,
//		}).First(&inv); res.RowsAffected == 0 {
//			tx.Rollback() //回滚之前的操作
//			return nil, status.Errorf(codes.NotFound, "该商品暂无库存信息")
//		}
//
//		if goodInfo.Num > inv.Stocks {
//			tx.Rollback()
//			return nil, status.Errorf(codes.ResourceExhausted, "库存不足")
//		}
//		inv.Stocks -= goodInfo.Num
//		tx.Save(&inv)
//	}
//	tx.Commit()
//	return &emptypb.Empty{}, nil
//
//}

//func (InventoryServer) Sell(ctx context.Context, req *proto.SellInfo) (*emptypb.Empty, error) { //乐观锁
//
//	tx := global.DB.Begin()
//
//	for _, goodInfo := range req.GoodsInfo {
//		for { //乐观锁 版本号冲突的时候的重试机制
//			var inv model.Inventory
//			if res := global.DB.Where(&model.Inventory{ //乐观锁
//				Goods: goodInfo.GoodsId,
//			}).First(&inv); res.RowsAffected == 0 {
//				tx.Rollback() //回滚之前的操作
//				return nil, status.Errorf(codes.NotFound, "该商品暂无库存信息")
//			}
//
//			if goodInfo.Num > inv.Stocks {
//				tx.Rollback()
//				return nil, status.Errorf(codes.ResourceExhausted, "库存不足")
//			}
//			inv.Stocks -= goodInfo.Num
//			//gorm默认不更新零值 得加一个.select来声明强制更新
//			if updateRes := global.DB.Model(&model.Inventory{}).Select("Stocks", "version").Where("goods = ? and version = ?", goodInfo.GoodsId, inv.Version).Updates(model.Inventory{ //updates更新多个字段
//				Stocks:  inv.Stocks,
//				Version: inv.Version + 1,
//			}); updateRes.RowsAffected == 0 {
//				fmt.Println("update冲突")
//				continue
//			}
//			break
//		}
//
//		//tx.Save(&inv)
//	}
//
//	tx.Commit()
//	return &emptypb.Empty{}, nil
//
//}

//redis的上锁方式
func (InventoryServer) Sell(ctx context.Context, req *proto.SellInfo) (*emptypb.Empty, error) {

	tx := global.DB.Begin()
	for _, goodInfo := range req.GoodsInfo {
		mutex := global.RS.NewMutex(fmt.Sprintf("goods_%d", goodInfo.GoodsId)) //建一个key 如果这个key存在则获取它(setnx)
		fmt.Println(mutex)
		if err := mutex.Lock(); err != nil { //可能会超时
			return nil, status.Errorf(codes.Internal, "redis获取锁失败")
		}
		var inv model.Inventory
		if res := global.DB.Where(&model.Inventory{
			Goods: goodInfo.GoodsId,
		}).First(&inv); res.RowsAffected == 0 {
			tx.Rollback() //回滚之前的操作
			return nil, status.Errorf(codes.NotFound, "该商品暂无库存信息")
		}

		if goodInfo.Num > inv.Stocks {
			tx.Rollback()
			return nil, status.Errorf(codes.ResourceExhausted, "库存不足")
		}

		inv.Stocks -= goodInfo.Num

		tx.Save(&inv)
		if ok, err := mutex.Unlock(); !ok || err != nil {
			return nil, status.Errorf(codes.Internal, "redis释放锁失败")
		}
	}

	tx.Commit()
	return &emptypb.Empty{}, nil

}

//库存归还 1.订单超时情况下 2.订单创建失败，归还之前扣减的库存
func (InventoryServer) Reback(ctx context.Context, req *proto.SellInfo) (*emptypb.Empty, error) {
	//同一个订单 多个商品 这些商品有问题 也得要回归库存
	tx := global.DB.Begin()
	for _, goodInfo := range req.GoodsInfo {
		var inv model.Inventory
		if res := global.DB.Where(&model.Inventory{
			Goods: goodInfo.GoodsId,
		}).First(&inv); res.RowsAffected == 0 {
			tx.Rollback() //回滚之前的操作
			return nil, status.Errorf(codes.NotFound, "该商品暂无库存信息")
		}

		inv.Stocks += goodInfo.Num
		tx.Save(&inv)
	}
	tx.Commit()
	return &emptypb.Empty{}, nil
}
