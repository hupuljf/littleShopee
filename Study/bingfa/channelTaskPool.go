package bingfa

//这个任务池维护的是一个令牌池 不是缓冲的go routine 只有拿到令牌才起一个协程
type TaskPool struct {
	ch chan struct{}
}

func NewTask(limit int) *TaskPool {
	t := &TaskPool{
		ch: make(chan struct{}, limit),
	}
	for i := 0; i < limit; i++ {
		t.ch <- struct{}{} //空结构体不占空间
	}
	return t
}

func (t *TaskPool) Do(f func()) {
	token := <-t.ch //没拿到令牌前都是阻塞状态
	//异步执行
	go func() {
		f()
		t.ch <- token //任务执行完以后把令牌还回去
	}()

	//同步执行
	//f()
	//t.ch <- token //任务执行完以后把令牌还回去

}
