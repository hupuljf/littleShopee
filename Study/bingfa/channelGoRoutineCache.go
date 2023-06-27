package bingfa

//维护的是goroutine 队列
type TaskPoolWithCache struct {
	ch chan func()
}

func NewTaskWithCache(limit int) *TaskPoolWithCache {
	t := &TaskPoolWithCache{
		ch: make(chan func(), limit),
	}
	for i := 0; i < limit; i++ {
		go func() { //开了limit个协程 在等待有函数调用它的信号量
			for {
				select {
				case fun := <-t.ch: //接收到了信号量
					fun()
				}
			}
		}()
	}
	return t
}

func (t *TaskPoolWithCache) Do(f func()) {
	t.ch <- f
}
