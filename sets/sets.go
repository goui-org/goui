package sets

type Set[T comparable] struct {
	m map[T]struct{}
	s []T
}

func New[T comparable]() *Set[T] {
	return &Set[T]{
		m: make(map[T]struct{}),
	}
}

func (s *Set[T]) Add(value T) {
	if _, ok := s.m[value]; ok {
		return
	}
	s.m[value] = struct{}{}
	s.s = append(s.s, value)
}

func (s *Set[T]) Has(value T) bool {
	_, ok := s.m[value]
	return ok
}

func (s *Set[T]) Delete(value T) {
	if !s.Has(value) {
		return
	}
	delete(s.m, value)
	for i := 0; i < len(s.s); i++ {
		if s.s[i] == value {
			s.s = append(s.s[:i], s.s[i+1:]...)
			return
		}
	}
}

func (s *Set[T]) Slice() []T {
	resp := make([]T, len(s.s))
	copy(resp, s.s)
	return resp
}
