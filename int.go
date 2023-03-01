package util

type Num interface {
	~int | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~int8 | ~int16 | ~int32 | ~int64
}

func Max[T Num](a T, b T) T {
	if a > b {
		return a
	}
	return b
}

func Min[T Num](a, b T) T {
	if a < b {
		return a
	}
	return b
}
