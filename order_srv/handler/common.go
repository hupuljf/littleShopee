package handler

import (
	"fmt"
	"math/rand"
	"time"
)

func GenerateOrderSn(userId int32) string {
	//时间加用户id加随机数
	now := time.Now()
	rand.Seed(time.Now().UnixNano()) //根据时间来 保证每次的随机数种子都是不一样的
	orderSn := fmt.Sprintf("%d%d%d%d%d%d%d%d", now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Nanosecond(), userId, rand.Intn(90)+10)
	return orderSn
}
