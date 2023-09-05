package slicex

import (
	"golang.org/x/exp/slices"

	"github.com/tomatocuke/candy/typex"
)

// 元素位置 - 从前往后
func Index[E comparable](es []E, e E) int {
	return slices.Index(es, e)
}

// 元素位置 - 从后往前
func LastIndex[E comparable](es []E, e E) int {
	for i := len(es) - 1; i >= 0; i-- {
		if es[i] == e {
			return i
		}
	}
	return -1
}

// 是否包含元素
func Contains[E comparable](es []E, e E) bool {
	return slices.Contains(es, e)
}

// 和
func Sum[E typex.Integer](es []E) E {
	var result E
	for _, e := range es {
		result += e
	}
	return result
}

// 平均值
func Average[E typex.Integer](es []E) float64 {
	if len(es) == 0 {
		return 0
	}

	return float64(Sum(es)) / float64(len(es))
}

// 去重
func Distinct[E comparable](es []E) []E {
	result := make([]E, 0, len(es))
	for _, e := range es {
		if !Contains(result, e) {
			result = append(result, e)
		}
	}
	return result
}

// 过滤
func Filter[E interface{}](es []E, f func(E) bool) []E {
	result := make([]E, 0, len(es))
	for _, e := range es {
		if f(e) {
			result = append(result, e)
		}
	}
	return result
}

// 合并
func Merge[E interface{}](es ...[]E) []E {
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
