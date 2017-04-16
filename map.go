package concurrent

import (
	"sync"
)

const (
	ShardNum = 10
)

func getShardIndex(k HashKey) uint32 {
	return k.hash() % ShardNum
}

type HashKey interface {
	hash() uint32
}

type innerMap struct {
	m map[HashKey]interface{}
	sync.RWMutex
}

type ConcurrentMap []*innerMap

func New() ConcurrentMap {
	var cm ConcurrentMap

	for i := 0; i < ShardNum; i++ {
		cm = append(cm, &innerMap{
			m: make(map[HashKey]interface{}),
		})
	}

	return cm
}

func (cm ConcurrentMap) Len() int {
	l := 0
	for i := 0; i < ShardNum; i++ {
		cm[i].RLock()
		l += len(cm[i].m)
		cm[i].RUnlock()
	}

	return l
}

func (cm ConcurrentMap) Set(k HashKey, v interface{}) {
	i := getShardIndex(k)
	cm[i].Lock()
	cm[i].m[k] = v
	cm[i].Unlock()
}

func (cm ConcurrentMap) Get(k HashKey) interface{} {
	i := getShardIndex(k)
	cm[i].RLock()
	v := cm[i].m[k]
	cm[i].RUnlock()
	return v
}

func (cm ConcurrentMap) Has(k HashKey) bool {
	i := getShardIndex(k)
	cm[i].RLock()
	_, ok := cm[i].m[k]
	cm[i].RUnlock()
	return ok
}

func (cm ConcurrentMap) Items() (m map[HashKey]interface{}) {
	m = make(map[HashKey]interface{})
	for i := 0; i < ShardNum; i++ {
		cm[i].RLock()
		for k, v := range cm[i].m {
			m[k] = v
		}
		cm[i].RUnlock()
	}

	return
}

func (cm ConcurrentMap) PutIfAbsent(k HashKey, v interface{}) (interface{}, bool) {
	i := getShardIndex(k)
	cm[i].Lock()
	defer cm[i].Unlock()
	prev, ok := cm[i].m[k]
	if ok {
		return prev, false
	} else {
		cm[i].m[k] = v
		return v, true
	}
}

func (cm ConcurrentMap) Keys() (keys []HashKey) {
	for i := 0; i < ShardNum; i++ {
		cm[i].RLock()
		for k := range cm[i].m {
			keys = append(keys, k)
		}
		cm[i].RUnlock()
	}

	return
}

func (cm ConcurrentMap) Values() (values []interface{}) {
	for i := 0; i < ShardNum; i++ {
		cm[i].RLock()
		for _, v := range cm[i].m {
			values = append(values, v)
		}
		cm[i].RUnlock()
	}

	return
}

func (cm ConcurrentMap) Clear() {
	for i := 0; i < ShardNum; i++ {
		cm[i].RLock()
		cm[i].m = make(map[HashKey]interface{})
		cm[i].RUnlock()
	}
}

func (cm ConcurrentMap) Update(m map[HashKey]interface{}) {
	for k, v := range m {
		i := getShardIndex(k)
		cm[i].Lock()
		cm[i].m[k] = v
		cm[i].Unlock()
	}
}

func (cm ConcurrentMap) IsEmpty() bool {
	return cm.Len() == 0
}
