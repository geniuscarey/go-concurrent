package concurrent

import (
	"sync"
)

const (
	ShardNum = 10
)

func getShardIndex(k HashKey) uint32 {
	return k.Hash() % ShardNum
}

type HashKey interface {
	Hash() uint32
}

type innerMap struct {
	m map[HashKey]interface{}
	sync.RWMutex
}

type Map []*innerMap

func NewMap() Map {
	var cm Map

	for i := 0; i < ShardNum; i++ {
		cm = append(cm, &innerMap{
			m: make(map[HashKey]interface{}),
		})
	}

	return cm
}

func (cm Map) Len() int {
	l := 0
	for i := 0; i < ShardNum; i++ {
		cm[i].RLock()
		l += len(cm[i].m)
		cm[i].RUnlock()
	}

	return l
}

func (cm Map) Set(k HashKey, v interface{}) {
	i := getShardIndex(k)
	cm[i].Lock()
	cm[i].m[k] = v
	cm[i].Unlock()
}

func (cm Map) Get(k HashKey) interface{} {
	i := getShardIndex(k)
	cm[i].RLock()
	v := cm[i].m[k]
	cm[i].RUnlock()
	return v
}

func (cm Map) Delete(k HashKey) interface{} {
	i := getShardIndex(k)
	cm[i].RLock()
	v := cm[i].m[k]
	delete(cm[i].m, k)
	cm[i].RUnlock()
	return v
}

func (cm Map) Has(k HashKey) bool {
	i := getShardIndex(k)
	cm[i].RLock()
	_, ok := cm[i].m[k]
	cm[i].RUnlock()
	return ok
}

func (cm Map) Items() (m map[HashKey]interface{}) {
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

func (cm Map) PutIfAbsent(k HashKey, v interface{}) (interface{}, bool) {
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

func (cm Map) Keys() (keys []HashKey) {
	for i := 0; i < ShardNum; i++ {
		cm[i].RLock()
		for k := range cm[i].m {
			keys = append(keys, k)
		}
		cm[i].RUnlock()
	}

	return
}

func (cm Map) Values() (values []interface{}) {
	for i := 0; i < ShardNum; i++ {
		cm[i].RLock()
		for _, v := range cm[i].m {
			values = append(values, v)
		}
		cm[i].RUnlock()
	}

	return
}

func (cm Map) Clear() {
	for i := 0; i < ShardNum; i++ {
		cm[i].RLock()
		cm[i].m = make(map[HashKey]interface{})
		cm[i].RUnlock()
	}
}

func (cm Map) Update(m map[HashKey]interface{}) {
	for k, v := range m {
		i := getShardIndex(k)
		cm[i].Lock()
		cm[i].m[k] = v
		cm[i].Unlock()
	}
}

func (cm Map) IsEmpty() bool {
	return cm.Len() == 0
}
