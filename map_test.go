package concurrent

import (
	"reflect"
	"testing"
	"time"
)

type NodeId uint32

func (id NodeId) Hash() uint32 {
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
	for _, _ = range m.Items() {
		//
	}
}

type timestamp time.Time

func (t timestamp) Compare(v interface{}) bool {
	o := v.(timestamp)
	return time.Time(t).Sub(time.Time(o)) > 0
}

func TestMapPutIfNewer(t *testing.T) {
	m := NewMap()
	prev := timestamp(time.Now())
	time.Sleep(time.Millisecond)
	m.Set(NodeId(3), timestamp(time.Now()))
	time.Sleep(time.Millisecond)
	newer := m.PutIfNewer(NodeId(3), timestamp(time.Now()))
	if !newer {
		t.Error("PutIfNewer failed")
	}

	newer = m.PutIfNewer(NodeId(3), prev)
	if newer {
		t.Error("PutIfNewer failed")
	}
}

func TestMapSortedKeys(t *testing.T) {
	m := NewMap()
	m.Set(NodeId(3), 4)
	m.Set(NodeId(63), 5)
	m.Set(NodeId(42), 5)
	m.Set(NodeId(14), 5)
	if m.SortedKeys()[1] != NodeId(14) {
		t.Error("Sort is not correct")
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
