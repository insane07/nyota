package gods

import (
	"sync"
	"unsafe"
)

type StringPool struct {
	sync.RWMutex
	Limit int
	pool  map[string]string
	size  int64
}

func NewStringPool(limit int) *StringPool {
	return &StringPool{
		Limit: limit,
		pool:  make(map[string]string),
	}
}

func (sp *StringPool) Get(s string) string {
	sp.RLock()
	sobj, ok := sp.pool[s]
	sp.RUnlock()
	if ok {
		return sobj
	}

	sp.Lock()
	defer sp.Unlock()
	sobj, ok = sp.pool[s]
	if ok {
		return sobj
	}

	if sp.Limit > 0 && len(sp.pool) >= sp.Limit {
		return s
	}
	sp.pool[s] = s
	sp.size = int64(unsafe.Sizeof(s)) + int64(len(s))
	return s
}

func (sp *StringPool) Count() int {
	sp.RLock()
	defer sp.RUnlock()
	return len(sp.pool)
}

func (sp *StringPool) Size() int64 {
	sp.RLock()
	defer sp.RUnlock()
	return sp.size
}
