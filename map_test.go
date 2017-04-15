package concurrent

import "testing"
import "reflect"

func TestMapCreation(t *testing.T) {
	m := New()
	if m == nil {
		t.Error("map is null.")
	}

	if m.Len() != 0 {
		t.Error("empty map's length should be zero")
	}

}

type NodeId uint32

func (id NodeId) hash() uint32 {
	return uint32(id)
}

func TestMapSet(t *testing.T) {
	m := New()
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
	m := New()
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
	m := New()
	m.Set(NodeId(3), int(4))
	m.Set(NodeId(4), int(5))
	if len(m.Keys()) != 2 {
		t.Error("Return value error")
	}
}

func TestMapValues(t *testing.T) {
	m := New()
	m.Set(NodeId(3), int(4))
	m.Set(NodeId(4), int(5))
	if len(m.Values()) != 2 {
		t.Error("Return value error")
	}
}
