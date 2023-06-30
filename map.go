package candy

import (
	"sync"
	"time"
)

// ========================= 函数 =========================

func Keys[K comparable, V any](m map[K]V) []K {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func Values[K comparable, V any](m map[K]V) []V {
	values := make([]V, 0, len(m))
	for _, v := range m {
		values = append(values, v)
	}
	return values
}

// ========================= 带过期时间的Map =========================

type Map[K comparable, V any] struct {
	// 键值过期时间(秒)
	expireSecond uint32
	// 数据
	data sync.Map
}

type mapItem[V any] struct {
	Value    V
	ExpireAt uint32
}

func (m *Map[K, V]) SetExpires(d time.Duration) {
	m.expireSecond = uint32(d / time.Second)
}

// 存储。如果无过期时间，则直接存储元素。有过期时间，存储mapItem
func (m *Map[K, V]) Set(key K, value V) {
	var save any
	if m.expireSecond == 0 {
		save = value
	} else {
		save = mapItem[V]{
			ExpireAt: uint32(time.Now().Unix()) + m.expireSecond,
			Value:    value,
		}
	}
	m.data.Store(key, save)
}

func (m *Map[K, V]) Get(key K) (V, bool) {
	v, ok := m.data.Load(key)

	if ok {
		// 无过期设置
		if m.expireSecond == 0 {
			return v.(V), true
		}

		item := v.(mapItem[V])
		// 未过期
		if time.Now().Unix() < int64(item.ExpireAt) {
			return item.Value, true
		}

		// 过期删除
		m.data.Delete(key)
	}

	var empty V
	return empty, false
}

func (m *Map[K, V]) Has(key K) bool {
	_, ok := m.Get(key)
	return ok
}

func (m *Map[K, V]) Del(key K) {
	m.data.Delete(key)
}

// ========================= Set =========================

type Set[E comparable] struct {
	mu   sync.RWMutex
	data map[E]struct{}
}

func NewSet[E comparable]() *Set[E] {
	return &Set[E]{}
}

func (s *Set[E]) Set(es ...E) {
	s.mu.Lock()
	if s.data == nil {
		s.data = make(map[E]struct{})
	}
	for _, e := range es {
		s.data[e] = struct{}{}
	}
	s.mu.Unlock()
}

func (s *Set[E]) Del(es ...E) {
	s.mu.Lock()
	for _, e := range es {
		delete(s.data, e)
	}
	s.mu.Unlock()
}

func (s *Set[E]) Has(e E) bool {
	s.mu.RLock()
	_, ok := s.data[e]
	s.mu.RUnlock()
	return ok
}

func (s *Set[E]) Len() int {
	s.mu.RLock()
	l := len(s.data)
	s.mu.RUnlock()
	return l
}

func (s *Set[E]) Clear() {
	s.mu.Lock()
	s.data = make(map[E]struct{})
	s.mu.Unlock()
}
