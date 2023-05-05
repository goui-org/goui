package goui

import "sync"

type store[K comparable, V any] struct {
	m  map[K]V
	mu sync.RWMutex
}

func newStore[K comparable, V any]() *store[K, V] {
	return &store[K, V]{
		m: make(map[K]V),
	}
}

func (s *store[K, V]) get(k K) V {
	s.mu.RLock()
	v := s.m[k]
	s.mu.RUnlock()
	return v
}

func (s *store[K, V]) set(k K, v V) {
	s.mu.Lock()
	s.m[k] = v
	s.mu.Unlock()
}

func (s *store[K, V]) delete(k K) {
	s.mu.Lock()
	delete(s.m, k)
	s.mu.Unlock()
}

func (s *store[K, V]) clear() {
	s.mu.Lock()
	for k := range s.m {
		delete(s.m, k)
	}
	s.mu.Unlock()
}

func (s *store[K, V]) all() []V {
	s.mu.RLock()
	defer s.mu.RUnlock()
	vs := make([]V, 0, len(s.m))
	for _, v := range s.m {
		vs = append(vs, v)
	}
	return vs
}
