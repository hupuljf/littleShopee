package main

import (
	"crypto/sha512"
	"fmt"
	"github.com/anaskhan96/go-password-encoder"
	"strings"
)

//用于同步表结构 或者 做一些功能测试
//md5加密与盐值加密
//md5加密事实上不可逆 但是简单密码会被暴力破解（会有人维护彩虹表 直接查询彩虹表）
//基于以上问题 需要给密码加盐 加盐只的就是随机字符串加用户密码
func EncodePassword(rawPassword string) string {
	options := &password.Options{16, 100, 32, sha512.New}
	salt, encodedPwd := password.Encode(rawPassword, options)
	newPassword := fmt.Sprintf("pbkdf2-sha512$%s$%s", salt, encodedPwd)
	return newPassword
}
func DecodePassword(pwd string, encoded string) bool {
	options := &password.Options{16, 100, 32, sha512.New}
	passwordInfo := strings.Split(encoded, "$")
	return password.Verify(pwd, passwordInfo[1], passwordInfo[2], options)
}
func main() {
	//user服务的数据库是shop_usr_srv
	//dsn := "root:8971841xm@tcp(localhost:3306)/shop_usr_srv?charset=utf8mb4&parseTime=True&loc=Local"
	//newLogger := logger.New(
	//	log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
	//	logger.Config{
	//		SlowThreshold: time.Second, // 慢 SQL 阈值
	//		LogLevel:      logger.Info, // Log level
	//		Colorful:      true,        // 禁用彩色打印
	//	},
	//)
	//
	//// 全局模式
	//db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
	//	Logger: newLogger,
	//})
	//if err != nil {
	//	panic(err)
	//}
	//_ = db.AutoMigrate(&model.User{})

	//db.Create(&model.User{
	//	NickName: "Fuck",
	//})
	//db.First(&model.User{}, "nick_name = ?", "Fuck")

	//Using custom options
	//options := &password.Options{16, 100, 32, sha512.New}
	//salt, encodedPwd := password.Encode("mimashihehe1.14", options)
	//newPassword := fmt.Sprintf("$pbkdf2-sha512$%s$%s", salt, encodedPwd)
	//fmt.Println(len(newPassword))
	////存进数据库的是这个new password
	//fmt.Println(newPassword)
	//
	//passwordInfo := strings.Split(newPassword, "$")
	//fmt.Println(passwordInfo)
	//check := password.Verify("mimashihehe1.14", passwordInfo[2], passwordInfo[3], options)
	//fmt.Println(check) // true
	pwd := EncodePassword("mimimimi")
	fmt.Println(pwd)
	fmt.Println(DecodePassword("mimimimi", pwd))
}
