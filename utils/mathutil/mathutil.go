package mathutil

type sortable interface {
	int | uint | int8 | uint8 | int16 | uint16 | int32 | uint32 | int64 | uint64 | string
}

func Min[T sortable](a, b T) T {
	if a < b {
		return a
	}
	return b
}
