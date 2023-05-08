package concurrentmap

import "sync"

type Map[K comparable, V any] struct {
	m  map[K]V
	mu sync.RWMutex
}

func New[K comparable, V any]() *Map[K, V] {
	return &Map[K, V]{
		m: make(map[K]V),
	}
}

func (s *Map[K, V]) Get(k K) V {
	s.mu.RLock()
	v := s.m[k]
	s.mu.RUnlock()
	return v
}

func (s *Map[K, V]) Set(k K, v V) {
	s.mu.Lock()
	s.m[k] = v
	s.mu.Unlock()
}

func (s *Map[K, V]) Delete(k K) {
	s.mu.Lock()
	delete(s.m, k)
	s.mu.Unlock()
}

func (s *Map[K, V]) Clear() {
	s.mu.Lock()
	for k := range s.m {
		delete(s.m, k)
	}
	s.mu.Unlock()
}

func (s *Map[K, V]) Len() int {
	s.mu.RLock()
	l := len(s.m)
	s.mu.RUnlock()
	return l
}

func (s *Map[K, V]) AllValues() []V {
	s.mu.RLock()
	defer s.mu.RUnlock()
	vs := make([]V, 0, len(s.m))
	for _, v := range s.m {
		vs = append(vs, v)
	}
	return vs
}
