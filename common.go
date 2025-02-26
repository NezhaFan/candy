package candy

import (
	"fmt"
	"math"
	"reflect"
	"runtime"
	"strings"
	"unicode/utf8"
)

// 字符串长度
func StrLen(s string) int {
	return utf8.RuneCountInString(s)
}

// 除法。保留x位小数
func Div[T Number](a, b T, decimals uint8) float64 {
	if b == 0 {
		return 0
	}
	return Round(float64(a)/float64(b), decimals)
}

// 保留x位小数
func Round(n float64, decimals uint8) float64 {
	ratio := math.Pow(10, float64(decimals))
	return math.Round(n*ratio) / ratio
}

// 结构体、结构体指针、结构体切片
func GetTags(v any, tag string) []string {
	var tp reflect.Type
	if x, ok := v.(reflect.Type); ok {
		tp = x
	} else {
		tp = reflect.TypeOf(v)
	}

	if tp.Kind() == reflect.Ptr {
		tp = tp.Elem()
	}

	if tp.Kind() == reflect.Slice {
		tp = tp.Elem()
	}

	n := tp.NumField()
	result := make([]string, 0, n)
	for i := 0; i < n; i++ {
		field := tp.Field(i)
		s := field.Tag.Get(tag)
		if s == "-" {
			continue
		}
		// json:"xxx,omitempty"
		if s != "" {
			if idx := strings.IndexByte(s, ','); idx > -1 {
				s = s[0:idx]
			}
		}
		// 如果有标签则记录，如果没标签，判断是否为嵌套的结构体
		if s != "" {
			result = append(result, s)
		} else if field.Type.Kind() == reflect.Struct {
			result = append(result, GetTags(field.Type, tag)...)
		}
	}

	return result
}

// 追踪调用位置
func Callers() []string {
	pc := make([]uintptr, 8)
	n := runtime.Callers(3, pc)
	if n == 0 {
		return nil
	}

	callers := make([]string, 0, n)

	pc = pc[:n]
	for _, p := range pc {
		fn := runtime.FuncForPC(p)
		if fn == nil {
			continue
		}
		file, line := fn.FileLine(p)
		arr := strings.Split(file, "/")
		if len(arr) > 3 {
			file = strings.Join(arr[len(arr)-3:], "/")
		}
		callers = append(callers, fmt.Sprintf("%s:%d", file, line))
	}

	return callers
}
