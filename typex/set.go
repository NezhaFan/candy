package typex

import "sync"

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
