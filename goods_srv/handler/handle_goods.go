package handler

import (
	"context"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"littleShopee/goods_srv/global"
	"littleShopee/goods_srv/model"
	"littleShopee/goods_srv/proto"
)

func ModelToResponse(goods model.Goods) proto.GoodsInfoResponse {
	return proto.GoodsInfoResponse{
		Id:              goods.ID,
		CategoryId:      goods.CategoryID,
		Name:            goods.Name,
		GoodsSn:         goods.GoodsSn,
		ClickNum:        goods.ClickNum,
		SoldNum:         goods.SoldNum,
		FavNum:          goods.FavNum,
		MarketPrice:     goods.MarketPrice,
		ShopPrice:       goods.ShopPrice,
		GoodsBrief:      goods.GoodsBrief,
		ShipFree:        goods.ShipFree,
		GoodsFrontImage: goods.GoodsFrontImage,
		IsNew:           goods.IsNew,
		IsHot:           goods.IsHot,
		OnSale:          goods.OnSale,
		DescImages:      goods.DescImages,
		Images:          goods.Images,
		Category: &proto.CategoryBriefInfoResponse{
			Id:   goods.Category.ID,
			Name: goods.Category.Name, //不preload的话 实际上 goods这个表里面存的只有category_id
		},
		Brand: &proto.BrandInfoResponse{
			Id:   goods.Brands.ID,
			Name: goods.Brands.Name,
			Logo: goods.Brands.Logo,
		},
	}
}

//type GoodsFilterRequest struct {
//	state         protoimpl.MessageState
//	sizeCache     protoimpl.SizeCache
//	unknownFields protoimpl.UnknownFields
//
//	PriceMin    int32  `protobuf:"varint,1,opt,name=priceMin,proto3" json:"priceMin,omitempty"`
//	PriceMax    int32  `protobuf:"varint,2,opt,name=priceMax,proto3" json:"priceMax,omitempty"`
//	IsHot       bool   `protobuf:"varint,3,opt,name=isHot,proto3" json:"isHot,omitempty"`
//	IsNew       bool   `protobuf:"varint,4,opt,name=isNew,proto3" json:"isNew,omitempty"`
//	IsTab       bool   `protobuf:"varint,5,opt,name=isTab,proto3" json:"isTab,omitempty"`
//	TopCategory int32  `protobuf:"varint,6,opt,name=topCategory,proto3" json:"topCategory,omitempty"`
//	Pages       int32  `protobuf:"varint,7,opt,name=pages,proto3" json:"pages,omitempty"`
//	PagePerNums int32  `protobuf:"varint,8,opt,name=pagePerNums,proto3" json:"pagePerNums,omitempty"`
//	KeyWords    string `protobuf:"bytes,9,opt,name=keyWords,proto3" json:"keyWords,omitempty"`
//	Brand       int32  `protobuf:"varint,10,opt,name=brand,proto3" json:"brand,omitempty"`
//}

func (s GoodsServer) GoodsList(ctx context.Context, req *proto.GoodsFilterRequest) (*proto.GoodsListResponse, error) {
	/*
		1.判断各个参数是不是非空值
		2.如果为非空值 修改localdb
		3.最终查询结果是综合各个查询的local db的查询结果
	*/
	goodListResp := &proto.GoodsListResponse{}
	var goodsList []model.Goods
	localDB := global.DB.Model(model.Goods{}) //代表用这个db只查goods表
	if req.PriceMin > 0 {
		localDB = localDB.Where("shop_price > ?", req.PriceMin)
	}
	if req.PriceMax > 0 {
		localDB = localDB.Where("shop_price < ?", req.PriceMax)
	}
	if req.IsHot == true {
		localDB = localDB.Where(model.Goods{IsHot: true})
	}
	if req.IsNew == true {
		localDB = localDB.Where(model.Goods{IsNew: true})
	}
	if req.KeyWords != "" {
		localDB = localDB.Where("name LIKE ?", "%"+req.KeyWords+"%")
	}
	if req.Brand != 0 {
		localDB = localDB.Where("brand_id = ?", req.Brand)
	}
	var subQuery string
	if req.TopCategory != 0 {
		var category model.Category
		res := global.DB.First(&category, req.TopCategory)
		if res.RowsAffected == 0 {
			return nil, status.Errorf(codes.NotFound, "商品不存在")
		}
		if category.Level == 1 {
			subQuery = fmt.Sprintf("select id form category where parent_category_id in (select id from category where parent_category_id = %d)", req.TopCategory)
		} else if category.Level == 2 {
			subQuery = fmt.Sprintf("select id form category where parent_category_id = %d)", req.TopCategory)
		} else if category.Level == 3 {
			subQuery = fmt.Sprintf("%d", req.TopCategory)
		}
		localDB = localDB.Where(fmt.Sprintf("category_id in (%s)", subQuery))

	}
	var cnt int64
	localDB.Count(&cnt)
	goodListResp.Total = int32(cnt)
	localDB.Scopes(Paginate(int(req.Pages), int(req.PagePerNums))).Preload("Category").Preload("Brands").Find(&goodsList)

	//dataList := make([]*proto.GoodsInfoResponse, 0)
	for _, goodsInfo := range goodsList {
		gg := ModelToResponse(goodsInfo)
		goodListResp.Data = append(goodListResp.Data, &gg)

	}

	return goodListResp, nil

}

//现在用户提交订单有多个商品，你得批量查询商品的信息吧
func (s *GoodsServer) BatchGetGoods(ctx context.Context, req *proto.BatchGoodsIdInfo) (*proto.GoodsListResponse, error) {
	goodsListResponse := &proto.GoodsListResponse{}
	var goods []model.Goods

	//调用where并不会真正执行sql 只是用来生成sql的 当调用find， first才会去执行sql，
	result := global.DB.Where(req.Id).Find(&goods)
	for _, good := range goods {
		goodsInfoResponse := ModelToResponse(good)
		goodsListResponse.Data = append(goodsListResponse.Data, &goodsInfoResponse)
	}
	goodsListResponse.Total = int32(result.RowsAffected)
	return goodsListResponse, nil
}
func (s *GoodsServer) GetGoodsDetail(ctx context.Context, req *proto.GoodInfoRequest) (*proto.GoodsInfoResponse, error) {
	var goods model.Goods

	if result := global.DB.Preload("Category").Preload("Brands").First(&goods, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "商品不存在")
	}
	goodsInfoResponse := ModelToResponse(goods)
	return &goodsInfoResponse, nil
}

func (s *GoodsServer) CreateGoods(ctx context.Context, req *proto.CreateGoodsInfo) (*proto.GoodsInfoResponse, error) {
	var category model.Category
	if result := global.DB.First(&category, req.CategoryId); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "商品分类不存在")
	}

	var brand model.Brands
	if result := global.DB.First(&brand, req.BrandId); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "品牌不存在")
	}
	//先检查redis中是否有这个token
	//防止同一个token的数据重复插入到数据库中，如果redis中没有这个token则放入redis
	//这里没有看到图片文件是如何上传， 在微服务中 普通的文件上传已经不再使用
	goods := model.Goods{
		Brands:          brand,
		BrandsID:        brand.ID,
		Category:        category,
		CategoryID:      category.ID,
		Name:            req.Name,
		GoodsSn:         req.GoodsSn,
		MarketPrice:     req.MarketPrice,
		ShopPrice:       req.ShopPrice,
		GoodsBrief:      req.GoodsBrief,
		ShipFree:        req.ShipFree,
		Images:          req.Images,
		DescImages:      req.DescImages,
		GoodsFrontImage: req.GoodsFrontImage,
		IsNew:           req.IsNew,
		IsHot:           req.IsHot,
		OnSale:          req.OnSale,
	}

	//srv之间互相调用了
	tx := global.DB.Begin() //为什么要用事务 失败了可以回滚
	result := tx.Save(&goods)
	if result.Error != nil {
		tx.Rollback()
		return nil, result.Error
	}
	tx.Commit()
	return &proto.GoodsInfoResponse{
		Id: goods.ID,
	}, nil
}

func (s *GoodsServer) DeleteGoods(ctx context.Context, req *proto.DeleteGoodsInfo) (*emptypb.Empty, error) {
	if result := global.DB.Delete(&model.Goods{BaseModel: model.BaseModel{ID: req.Id}}, req.Id); result.Error != nil {
		return nil, status.Errorf(codes.NotFound, "商品不存在")
	}
	return &emptypb.Empty{}, nil
}

func (s *GoodsServer) UpdateGoods(ctx context.Context, req *proto.CreateGoodsInfo) (*emptypb.Empty, error) {
	var goods model.Goods

	if result := global.DB.First(&goods, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "商品不存在")
	}

	var category model.Category
	if result := global.DB.First(&category, req.CategoryId); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "商品分类不存在")
	}

	var brand model.Brands
	if result := global.DB.First(&brand, req.BrandId); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "品牌不存在")
	}

	goods.Brands = brand
	goods.BrandsID = brand.ID
	goods.Category = category
	goods.CategoryID = category.ID
	goods.Name = req.Name
	goods.GoodsSn = req.GoodsSn
	goods.MarketPrice = req.MarketPrice
	goods.ShopPrice = req.ShopPrice
	goods.GoodsBrief = req.GoodsBrief
	goods.ShipFree = req.ShipFree
	goods.Images = req.Images
	goods.DescImages = req.DescImages
	goods.GoodsFrontImage = req.GoodsFrontImage
	goods.IsNew = req.IsNew
	goods.IsHot = req.IsHot
	goods.OnSale = req.OnSale

	tx := global.DB.Begin()
	result := tx.Save(&goods)
	if result.Error != nil {
		tx.Rollback()
		return nil, result.Error
	}
	tx.Commit()
	return &emptypb.Empty{}, nil
}
