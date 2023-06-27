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

//泛型的使用 k 是可比较类型 V 是任意一个类型
type SafeMap[K comparable, V any] struct {
	values map[K]V
	lock   sync.RWMutex
}

//这个k已经存在的话 返回对应的值 loaded=true
//不存在则存 loaded=false
func (s *SafeMap[K, V]) LoadOrStore(k K, v V) (V, bool) {
	s.lock.RLock()
	oldVal, ok := s.values[k]
	s.lock.RUnlock() //不能用defer defer是return后执行的 如果没有return 后面获取不到写锁
	if ok {
		return oldVal, true
	}
	s.lock.Lock()
	oldVal, ok = s.values[k] //double check
	if ok {
		return oldVal, true
	}
	defer s.lock.Unlock()
	s.values[k] = v
	return v, false
}
