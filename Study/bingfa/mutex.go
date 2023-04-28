package bingfa

import "sync"

//把锁与资源封装在一起
type safeResource struct {
	resource map[string]string
	lock     sync.RWMutex
}

func (s *safeResource) Add(key, val string) {
	//想要操作资源先获得锁 先上锁
	s.lock.Lock()         //写锁 其他资源此时不能读不能写
	defer s.lock.Unlock() //释放写锁
	s.resource[key] = val
}

func (s *safeResource) Get(key string) string {
	s.lock.RLock()         //上读锁 其他资源可以读 不能写
	defer s.lock.RUnlock() //释放读锁 其他资源此时可以写了
	if val, ok := s.resource[key]; ok {
		return val
	} else {
		return "NAN"
	}
}
