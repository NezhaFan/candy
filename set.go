package candy

import (
	"sync"
)

type Set[E comparable] struct {
	mu   sync.RWMutex
	data map[E]struct{}
}

func (s *Set[E]) Set(es ...E) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.data == nil {
		s.data = make(map[E]struct{})
	}
	for _, e := range es {
		s.data[e] = struct{}{}
	}
}

func (s *Set[E]) Del(es ...E) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, e := range es {
		delete(s.data, e)
	}
}

func (s *Set[E]) Has(e E) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, ok := s.data[e]
	return ok
}

func (s *Set[E]) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.data)
}

func (s *Set[E]) Keys() []E {
	s.mu.RLock()
	defer s.mu.RUnlock()

	keys := make([]E, 0, len(s.data))
	for k := range s.data {
		keys = append(keys, k)
	}

	return keys
}
