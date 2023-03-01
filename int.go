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

func PadEnd(source string, size int, padStr string) string {
	return padAtPosition(source, size, padStr, 2)
}

func padAtPosition(str string, length int, padStr string, position int) string {
	if len(str) >= length {
		return str
	}

	if padStr == "" {
		padStr = " "
	}

	length = length - len(str)
	startPadLen := 0
	if position == 0 {
		startPadLen = length / 2
	} else if position == 1 {
		startPadLen = length
	}
	endPadLen := length - startPadLen

	charLen := len(padStr)
	leftPad := ""
	cur := 0
	for cur < startPadLen {
		leftPad += string(padStr[cur%charLen])
		cur++
	}

	cur = 0
	rightPad := ""
	for cur < endPadLen {
		rightPad += string(padStr[cur%charLen])
		cur++
	}

	return leftPad + str + rightPad
}

func Grow[S ~[]E, E any](s S, n int) S {
	if n < 0 {
		panic("cannot be negative")
	}
	if n -= cap(s) - len(s); n > 0 {
		// fmt.Println("当前n", n)
		s = append(s[:cap(s)], make([]E, n)...)[:len(s)]
	}
	return s
}

func Delete[S ~[]E, E any](s S, i, j int) S {
	_ = s[i:j] // bounds check

	return append(s[:i], s[j:]...)
}

func Clone[S ~[]E, E any](s S) S {

	// Preserve nil in case it matters.
	if s == nil {
		return nil
	}
	return append([]E{}, s...)
}

func Grow2[S ~[]E, E any](s S, n int) S {
	if n < 0 {
		panic("cannot be negative")
	}
	//
	// if len(s) == 0 {
	// 	return make(S, n)
	// }
	if n -= cap(s) - len(s); n > 0 {
		s = append(s[:cap(s)], make(S, n)...)[:len(s)]
		// s = append(s[:cap(s)], make([]E, n)...)[:len(s)]
	}
	return s
}
