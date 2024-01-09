package goui

// export let areDepsEqual = (a: any[], b: any[]): boolean => a.length === b.length && a.every((x, i) => a[i] === b[i]);

func areDepsEqual(a Deps, b Deps) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// Map creates a new slice with items that are mapped to new values
// according to the given m function. The given slice is not
// changed.
func Map[T any](s []T, m func(item T) *Elem) []*Elem {
	k := make([]*Elem, len(s))
	for i := range s {
		k[i] = m(s[i])
	}
	return k
}
