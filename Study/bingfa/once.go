package bingfa

import (
	"fmt"
	"sync"
)

//sync.once 确保某个动作最多执行一遍 初始化资源 单例模式

//1.单例模式 保证某个结构体只有一个实例

type Single struct {
	val string
}

func (s *Single) PrintSingle() {
	fmt.Println(s.val)
}

var Instance *Single
var OnceInit sync.Once

func InitSingleInstance() *Single {
	OnceInit.Do(func() {
		Instance = &Single{
			val: "test",
		}
		fmt.Println("Init")
	})
	return Instance
}

type Once struct {
	o sync.Once
}

func (pt *Once) TestPoint() { //指针的话调用100次都执行一次
	pt.o.Do(func() {
		fmt.Println("test")
	})
}

func (st Once) TestNOPoint() { //非指针 每次都再起一个sync包 100次就是100次
	st.o.Do(func() {
		fmt.Println("test NO Point")
	})
}
