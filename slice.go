package util

// 查找元素 请直接使用官方slices包

// 拼接切片 (一次分配内存，在多个大切片合并时效果好)
func SliceConcat[T any](slice []T, slices ...[]T) []T {
	size := len(slice)
	for _, s := range slices {
		size += len(s)
	}

	result := make([]T, size)
	i := copy(result, slice)

	for _, v := range slices {
		s := result[i:]
		i += copy(s, v)
	}

	return result
}

// func Values[KT comparable, VT any](m map[KT]VT) []T {
// 	r := make([]T, 0, len(m))
// 	for _, v := range m {
// 		r = append(r, v)
// 	}
// 	return r
// }
