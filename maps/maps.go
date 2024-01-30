package maps

type Map[K comparable, V any] struct {
	m map[K]V
	s []K
}

func New[K comparable, V any]() *Map[K, V] {
	return &Map[K, V]{
		m: make(map[K]V),
	}
}

func (m *Map[K, V]) Set(key K, value V) {
	if _, ok := m.m[key]; ok {
		m.m[key] = value
		return
	}
	m.m[key] = value
	m.s = append(m.s, key)
}

func (m *Map[K, V]) Has(key K) bool {
	_, ok := m.m[key]
	return ok
}

func (m *Map[K, V]) Get(key K) (V, bool) {
	v, ok := m.m[key]
	return v, ok
}

func (m *Map[K, V]) Delete(key K) {
	if !m.Has(key) {
		return
	}
	delete(m.m, key)
	for i := 0; i < len(m.s); i++ {
		if m.s[i] == key {
			m.s = append(m.s[:i], m.s[i+1:]...)
			return
		}
	}
}

type Entry[K comparable, V any] struct {
	Key   K
	Value V
}

func (m *Map[K, V]) Slice() []Entry[K, V] {
	resp := make([]Entry[K, V], 0, len(m.s))
	for k, v := range m.m {
		resp = append(resp, Entry[K, V]{
			Key:   k,
			Value: v,
		})
	}
	return resp
}
