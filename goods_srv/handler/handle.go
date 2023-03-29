package handler

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
	"littleShopee/goods_srv/global"
	"littleShopee/goods_srv/model"

	//"google.golang.org/protobuf/runtime/protoimpl"
	"littleShopee/goods_srv/proto"
)

type GoodsServer struct {
	proto.UnimplementedGoodsServer
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

//获取品牌列表
func (GoodsServer) BrandList(ctx context.Context, req *proto.BrandFilterRequest) (*proto.BrandListResponse, error) {
	var brands []model.Brands
	brandResp := new(proto.BrandListResponse)
	result := global.DB.Find(&brands) //所有brands 以此来获取brand的总数
	if result.Error != nil {          //代表查不到
		fmt.Println("错在哪儿", result.Error)
		return nil, result.Error
	}
	brandResp.Total = int32(result.RowsAffected)
	global.DB.Scopes(Paginate(int(req.Pages), int(req.PagePerNums))).Find(&brands)
	//brandResp.Data = append(brandResp.Data, append([]model.Brands{},brands...))
	for i := 0; i < len(brands); i++ {
		brandResp.Data = append(brandResp.Data, &proto.BrandInfoResponse{
			Id:   brands[i].ID,
			Name: brands[i].Name,
			Logo: brands[i].Logo,
		})
	}
	return brandResp, nil
}

func (GoodsServer) CreateBrand(ctx context.Context, req *proto.BrandRequest) (*proto.BrandInfoResponse, error) {
	//先判断品牌名是否创建过 创建过就报错
	if res := global.DB.Where("name = ?", req.Name).First(&model.Brands{}); res.RowsAffected == 1 {
		zap.S().Infof("crete brand fail:already existed")
		return nil, res.Error
	}
	newBrand := model.Brands{
		Name: req.Name,
		Logo: req.Logo,
	}
	global.DB.Save(&newBrand)
	return &proto.BrandInfoResponse{
		Id:   newBrand.ID,
		Name: newBrand.Name,
		Logo: newBrand.Logo,
	}, nil
}
func (GoodsServer) DeleteBrand(ctx context.Context, req *proto.BrandRequest) (*emptypb.Empty, error) {
	//根据id delete
	if res := global.DB.First(&model.Brands{}, req.Id); res.RowsAffected == 0 {
		zap.S().Infof("delete brand fail: no this id")
		return nil, res.Error
	}
	global.DB.Delete(&model.Brands{}, req.Id)
	return &emptypb.Empty{}, nil

}
func (GoodsServer) UpdateBrand(ctx context.Context, req *proto.BrandRequest) (*emptypb.Empty, error) {
	if res := global.DB.First(&model.Brands{}, req.Id); res.RowsAffected == 0 {
		zap.S().Infof("delete brand fail: no this id")
		return nil, res.Error
	}
	brands := model.Brands{}
	if req.Name != "" {
		brands.Name = req.Name
	}
	if req.Logo != "" {
		brands.Logo = req.Logo
	}
	global.DB.Save(&brands)
	return &emptypb.Empty{}, nil
}

//轮播图
func (s *GoodsServer) BannerList(ctx context.Context, req *emptypb.Empty) (*proto.BannerListResponse, error) {
	bannerListResponse := proto.BannerListResponse{}

	var banners []model.Banner
	result := global.DB.Find(&banners)
	bannerListResponse.Total = int32(result.RowsAffected)

	var bannerReponses []*proto.BannerResponse
	for _, banner := range banners {
		bannerReponses = append(bannerReponses, &proto.BannerResponse{
			Id:    banner.ID,
			Image: banner.Image,
			Index: banner.Index,
			Url:   banner.Url,
		})
	}

	bannerListResponse.Data = bannerReponses

	return &bannerListResponse, nil
}

func (GoodsServer) CreateBanner(ctx context.Context, req *proto.BannerRequest) (*proto.BannerResponse, error) {
	banner := model.Banner{} //图片重复无伤大雅

	banner.Image = req.Image
	banner.Index = req.Index
	banner.Url = req.Url

	global.DB.Save(&banner)

	return &proto.BannerResponse{Id: banner.ID}, nil
}

func (GoodsServer) DeleteBanner(ctx context.Context, req *proto.BannerRequest) (*emptypb.Empty, error) {
	if result := global.DB.Delete(&model.Banner{}, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "轮播图不存在")
	}
	return &emptypb.Empty{}, nil
}

func (GoodsServer) UpdateBanner(ctx context.Context, req *proto.BannerRequest) (*emptypb.Empty, error) {
	var banner model.Banner

	if result := global.DB.First(&banner, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "轮播图不存在")
	}

	if req.Url != "" {
		banner.Url = req.Url
	}
	if req.Image != "" {
		banner.Image = req.Image
	}
	if req.Index != 0 {
		banner.Index = req.Index
	}

	global.DB.Save(&banner)

	return &emptypb.Empty{}, nil
}
