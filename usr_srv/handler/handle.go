package handler

import (
	"context"
	"crypto/sha512"
	"fmt"
	"github.com/anaskhan96/go-password-encoder"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	"littleShopee/usr_srv/global"
	"littleShopee/usr_srv/model"
	"littleShopee/usr_srv/proto"
	"strings"
	"time"
)

//处理user的增删改查
type UserServer struct {
}

func Paginate(page, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page == 0 {
			page = 1
		}

		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

//类型转换
func ModelToRsponse(user model.User) proto.UserInfoResponse {
	//在grpc的message中字段有默认值，你不能随便赋值nil进去，容易出错
	//这里要搞清， 哪些字段是有默认值
	userInfoRsp := proto.UserInfoResponse{
		Id:       user.ID,
		PassWord: user.Password,
		NickName: user.NickName,
		Gender:   user.Gender,
		Role:     int32(user.Role),
		Mobile:   user.Mobile,
	}
	if user.Birthday != nil {
		userInfoRsp.BirthDay = uint64(user.Birthday.Unix())
	}
	return userInfoRsp
}

//返回对应第几页的user列表
func (server *UserServer) GetUserList(ctx context.Context, req *proto.PageInfo) (*proto.UserListResponse, error) {
	var users []model.User
	//users := make([]model.User, 0)
	//拿到所有users列表
	result := global.DB.Find(&users)
	if result.Error != nil {
		fmt.Println("错在哪儿", result.Error)
		return nil, result.Error
	}
	userList := new(proto.UserListResponse)
	(*userList).Total = int32(result.RowsAffected)
	global.DB.Scopes(Paginate(int(req.Pn), int(req.PSize))).Find(&users)
	for _, user := range users {
		userInfoRsp := ModelToRsponse(user)
		(*userList).Data = append((*userList).Data, &userInfoRsp)
	}

	return userList, nil
}

//通过mobile查询user
func (s *UserServer) GetUserByMobile(ctx context.Context, mobileReq *proto.MobileRequest) (*proto.UserInfoResponse, error) {
	var user model.User
	// SELECT * FROM users WHERE mobile =  mobileReq.Mobil
	result := global.DB.Where(&model.User{Mobile: mobileReq.Mobile}).First(&user)
	//res := global.DB.Where("mobile = ?",mobileReq.Mobile).First(&user)
	//res := global.DB.First(&user,"mobile = ?", mobileReq.Mobile)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "用户不存在")
	}
	if result.Error != nil {
		return nil, result.Error
	}

	userInfoRsp := ModelToRsponse(user)
	return &userInfoRsp, nil

}

func (s *UserServer) GetUserById(ctx context.Context, req *proto.IdRequest) (*proto.UserInfoResponse, error) {
	//通过id查询用户
	var user model.User
	result := global.DB.First(&user, req.Id)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "用户不存在")
	}
	if result.Error != nil {
		return nil, result.Error
	}

	userInfoRsp := ModelToRsponse(user)
	return &userInfoRsp, nil
}

func (s *UserServer) CreateUser(ctx context.Context, req *proto.CreateUserInfo) (*proto.UserInfoResponse, error) {
	//先通过unique的mobile查是否已经存在
	var user model.User
	res := global.DB.Where(&model.User{Mobile: req.Mobile}).First(&user)
	if res.RowsAffected > 0 {
		return nil, status.Errorf(codes.AlreadyExists, "用户已存在")
	}
	user.Mobile = req.Mobile
	user.NickName = req.NickName
	//密码加密
	options := &password.Options{16, 100, 32, sha512.New}
	salt, encodedPwd := password.Encode(req.PassWord, options)
	user.Password = fmt.Sprintf("pbkdf2-sha512$%s$%s", salt, encodedPwd)
	res = global.DB.Create(&user)
	if res.Error != nil {
		return nil, status.Errorf(codes.Internal, res.Error.Error())
	}
	userInfoRsp := ModelToRsponse(user)
	return &userInfoRsp, nil
}

func (s *UserServer) UpdateUser(ctx context.Context, req *proto.UpdateUserInfo) (*empty.Empty, error) {
	//个人中心更新用户
	var user model.User
	result := global.DB.First(&user, req.Id)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "用户不存在")
	}

	birthDay := time.Unix(int64(req.BirthDay), 0)
	user.NickName = req.NickName
	user.Birthday = &birthDay
	user.Gender = req.Gender

	result = global.DB.Save(&user)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}
	return &empty.Empty{}, nil
}

func (s *UserServer) CheckPassWord(ctx context.Context, req *proto.PasswordCheckInfo) (*proto.CheckResponse, error) {
	//校验密码
	options := &password.Options{16, 100, 32, sha512.New}
	passwordInfo := strings.Split(req.EncryptedPassword, "$")
	check := password.Verify(req.Password, passwordInfo[1], passwordInfo[2], options)
	return &proto.CheckResponse{Success: check}, nil
}
