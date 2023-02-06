package main

import (
	"container/heap"
	"fmt"
)

type element struct {
	value    int
	priority int
}
type priority_queue []element

func (pq priority_queue) Len() int {
	return len(pq)
}
func (pq priority_queue) Less(i, j int) bool {
	return pq[i].priority > pq[j].priority
}

//定义交换函数 交换的是啥 实际上并不需要修改队列本身
func (pq priority_queue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}
func (pq *priority_queue) Push(x interface{}) {
	//空接口 类型转换一下
	//*pq指的是指针指向的那个队列 如果是&pq的话就变成了指针的指针了
	*pq = append(*pq, x.(element))
}
func (pq *priority_queue) Pop() interface{} {
	res := (*pq)[len(*pq)-1]
	*pq = (*pq)[:len(*pq)-1]
	return res
}

func main() {
	pq := new(priority_queue)
	heap.Init(pq)
	heap.Push(pq, element{
		10,
		1,
	})
	heap.Push(pq, element{
		100,
		-100,
	})
	heap.Push(pq, element{
		1000,
		100,
	})
	fmt.Println(heap.Pop(pq))
	fmt.Println(heap.Pop(pq))
}
