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
	defer s.mu.RUnlock()
	return s.m[k]
}

func (s *store[K, V]) set(k K, v V) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.m[k] = v
}

func (s *store[K, V]) delete(k K) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.m, k)
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
