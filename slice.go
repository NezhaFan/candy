package candy

import "golang.org/x/exp/slices"

// 元素位置 - 从前往后
func IndexOfSlice[E comparable](es []E, e E) int {
	return slices.Index(es, e)
}

// 元素位置 - 从后往前
func LastIndexOfSlice[E comparable](es []E, e E) int {
	for i := len(es) - 1; i >= 0; i-- {
		if es[i] == e {
			return i
		}
	}
	return -1
}

// 是否包含元素
func ContainsOfSlice[E comparable](es []E, e E) bool {
	return slices.Contains(es, e)
}

// 和
func SumSlice[E Integer](es []E) E {
	var result E
	for _, e := range es {
		result += e
	}
	return result
}

// 平均值
func AverageOfSlice[E Integer](es []E) float64 {
	if len(es) == 0 {
		return 0
	}

	return float64(SumSlice(es)) / float64(len(es))
}

// 去重
func DistinctSlice[E comparable](es []E) []E {
	result := make([]E, 0, len(es))
	for _, e := range es {
		if !ContainsOfSlice(result, e) {
			result = append(result, e)
		}
	}
	return result
}

// 过滤
func FilterSlice[E any](es []E, f func(E) bool) []E {
	result := make([]E, 0, len(es))
	for _, e := range es {
		if f(e) {
			result = append(result, e)
		}
	}
	return result
}

// 拼接
func ConcatSlice[E any](es ...[]E) []E {
	size := 0
	for _, s := range es {
		size += len(s)
	}

	result := make([]E, size)
	for i, s := range es {
		i += copy(result[i:], s)
	}
	return result
}

// 转为any切片 (尽量不使用)
func ToAnySlice[E any](es []E) []any {
	result := make([]any, len(es))
	for i, e := range es {
		result[i] = e
	}
	return result
}
