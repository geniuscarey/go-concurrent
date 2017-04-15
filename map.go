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

type concurrentMap []*innerMap

func New() concurrentMap {
	var cm concurrentMap

	for i := 0; i < ShardNum; i++ {
		cm = append(cm, &innerMap{
			m: make(map[HashKey]interface{}),
		})
	}

	return cm
}

func (cm concurrentMap) Len() int {
	l := 0
	for i := 0; i < ShardNum; i++ {
		cm[i].RLock()
		l += len(cm[i].m)
		cm[i].RUnlock()
	}

	return l
}

func (cm concurrentMap) Set(k HashKey, v interface{}) {
	i := getShardIndex(k)
	cm[i].RLock()
	cm[i].m[k] = v
	cm[i].RUnlock()
}

func (cm concurrentMap) Get(k HashKey) interface{} {
	i := getShardIndex(k)
	cm[i].RLock()
	v := cm[i].m[k]
	cm[i].RUnlock()
	return v
}

func (cm concurrentMap) putIfAbsent(k HashKey, v interface{}) (interface{}, bool) {
	i := getShardIndex(k)
	cm[i].RLock()
	defer cm[i].RUnlock()
	prev, ok := cm[i].m[k]
	if ok {
		return prev, false
	} else {
		cm[i].m[k] = v
		return v, true
	}
}
