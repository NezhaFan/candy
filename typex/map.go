package typex

import "sync"

type Map[K comparable, V interface{}] struct {
	data sync.Map
}

func (m *Map[K, V]) Set(key K, value V) {
	m.data.Store(key, value)
}

func (m *Map[K, V]) Get(key K) (V, bool) {
	val, ok := m.data.Load(key)
	return val.(V), ok
}

func (m *Map[K, V]) MGet(keys []K) map[K]V {
	var empty V
	r := make(map[K]V, len(keys))
	for _, k := range keys {
		v, ok := m.Get(k)
		if ok {
			r[k] = v
		} else {
			r[k] = empty
		}
	}

	return r
}

func (m *Map[K, V]) Has(key K) bool {
	_, ok := m.data.Load(key)
	return ok
}

func (m *Map[K, V]) Del(key K) {
	m.data.Delete(key)
}
