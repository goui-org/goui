package goui

func areDepsEqual(a Deps, b Deps) bool {
	if len(a) == 0 || len(b) == 0 || len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func Map[T any](s []T, m func(item T) *Node) []*Node {
	k := make([]*Node, len(s))
	for i := range s {
		k[i] = m(s[i])
	}
	return k
}
