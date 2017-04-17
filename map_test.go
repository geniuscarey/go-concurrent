package concurrent

import (
	"fmt"
	"reflect"
	"testing"
)

type NodeId uint32

func (id NodeId) hash() uint32 {
	return uint32(id)
}

func TestMapCreation(t *testing.T) {
	m := NewMap()
	if m == nil {
		t.Error("map is null.")
	}

	if !m.IsEmpty() {
		t.Error("empty map's length should be zero")
	}
}

func TestMapSet(t *testing.T) {
	m := NewMap()
	m.Set(NodeId(3), int(4))

	v := m.Get(NodeId(3))
	if v != 4 {
		t.Error("Get value not eual with set value")
	}

	vi := v.(int)
	if reflect.TypeOf(vi).Name() != "int" {
		t.Error("Get value's type error")
	}
}

func TestMapPutIfAbsent(t *testing.T) {
	m := NewMap()
	m.Set(NodeId(3), int(4))

	prev, ok := m.PutIfAbsent(NodeId(3), 5)
	if ok {
		t.Error("Return status error")
	}

	if prev != 4 {
		t.Error("Return value error")
	}

	prev, ok = m.PutIfAbsent(NodeId(4), 5)

	if !ok {
		t.Error("Return status error")
	}

	if prev != 5 {
		t.Error("Return value error")
	}
}

func TestMapKeys(t *testing.T) {
	m := NewMap()
	m.Set(NodeId(3), int(4))
	m.Set(NodeId(4), int(5))
	if len(m.Keys()) != 2 {
		t.Error("Return value error")
	}
}

func TestMapValues(t *testing.T) {
	m := NewMap()
	m.Set(NodeId(3), int(4))
	m.Set(NodeId(4), int(5))
	if len(m.Values()) != 2 {
		t.Error("Return value error")
	}
}

func TestMapItems(t *testing.T) {
	m := NewMap()
	m.Set(NodeId(3), int(4))
	m.Set(NodeId(4), int(5))
	for k, v := range m.Items() {
		fmt.Println(k, v)
	}
}

func BenchmarkItems(b *testing.B) {
	m := NewMap()
	for i := 0; i < 10000; i++ {
		m.Set(NodeId(i), 0)
	}
	for n := 0; n < b.N; n++ {
		m.Items()
	}
}
