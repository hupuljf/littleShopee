package go_test_test

import (
	"fmt"
	"littleShopee/Study/bingfa"
	"sync"
	"testing"
	"time"
)

func TestDeferRLock(t *testing.T) {
	fmt.Println("test")
}

func TestOnceInstance(t *testing.T) {
	//go func() {
	//	s := bingfa.InitSingleInstance()
	//	fmt.Println("第一个")
	//	s.PrintSingle()
	//	fmt.Println(&s)
	//}()
	//time.Sleep(10 * time.Second)
	//go func() {
	//	s := bingfa.InitSingleInstance()
	//	fmt.Println("第二个")
	//	s.PrintSingle()`
	//	fmt.Println(&s)
	//}()
	bingfa.InitSingleInstance()
	bingfa.InitSingleInstance() // 只有一个打印 只调了一次

	pt := &bingfa.Once{}
	for i := 0; i < 10; i++ {
		pt.TestPoint()
	} //只执行了一遍

	st := bingfa.Once{}
	for i := 0; i < 10; i++ {
		st.TestNOPoint()
	}

}

func TestPool(t *testing.T) {
	pool := sync.Pool{
		New: func() any {
			fmt.Println("执行一次new")
			return 10000
		},
	}

	for i := 0; i < 10; i++ {
		val := pool.Get()
		//Get 会返回 Pool 已经存在的对象，如果没有，那么就走慢路径，也就是调用初始化的时候定义的 New 方法（也就是最开始定义的初始化行为）来初始化一个对象。
		fmt.Println(val)
		pool.Put(val) //put回去则不会再执行new了 pool作用不开很多变量 复用内存
	}
}

type user struct {
	uid int
	str string
}

func TestPoolCover(t *testing.T) {
	pool := sync.Pool{
		New: func() any {
			return &user{
				uid: 1,
				str: "1",
			}
		},
	}
	u1 := pool.Get().(*user)
	u1.uid = 100
	u1.str = "100"
	pool.Put(u1) //这个时候pool里面的对象已经被改了 //所以put前得先reset一下

	u2 := pool.Get().(*user)
	fmt.Println(u2)
}

func TestWaitGroup(t *testing.T) {
	wg := sync.WaitGroup{}
	for i := 1; i <= 10; i++ { //10个批次
		for j := 1 + (i-1)*10; j <= 10*i; j++ {
			wg.Add(1)
			x := j //注意得给j赋予一个新值 j都是同一个地址的 每个x:=都会赋予一个新的地址
			go func() {
				defer wg.Done()
				fmt.Println(x)
			}()
		}
		wg.Wait()
		time.Sleep(3 * time.Second)
	}
}

func TestChannel(t *testing.T) {
	ch := make(chan string, 4)
	go func() {
		str := <-ch
		fmt.Println("1" + str)
	}()
	go func() {
		str := <-ch
		fmt.Println("2" + str)
	}()
	go func() {
		str := <-ch
		fmt.Println("3" + str)
	}()
	ch <- "hello"
	ch <- "hello"
}

func TestChannelTask(t *testing.T) {
	taskPool := bingfa.NewTaskWithCache(2)
	taskPool.Do(func() {
		fmt.Println("task1 开始")

	})
	time.Sleep(time.Second)
	fmt.Println("task1 结束")
	taskPool.Do(func() {
		fmt.Println("task2 开始")

	})
	time.Sleep(time.Second)
	fmt.Println("task2 结束")
	taskPool.Do(func() {
		fmt.Println("task3 开始")

	})
	time.Sleep(time.Second)
	fmt.Println("task3 结束")

}
