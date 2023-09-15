package candy

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"net"
	"reflect"
	"strconv"
	"strings"
	"unicode/utf8"
	"unsafe"

	"github.com/tomatocuke/candy/typex"
)

// bytes转字符串
func BytesToString(b []byte) string {
	if len(b) == 0 {
		return ""
	}
	return unsafe.String(&b[0], len(b))
}

// 字符串转bytes
func StringToBytes(s string) (b []byte) {
	if s == "" {
		return nil
	}
	return unsafe.Slice(unsafe.StringData(s), len(s))
}

// 浮点数转字符串
func FloatToString(n float64) string {
	return strconv.FormatFloat(float64(n), 'f', -1, 64)
}

// 转字符串。优先根据类型来，能转json转json，最终才用fmt格式化
func AnyToString(x any) string {
	switch v := x.(type) {
	case string:
		return v
	case int:
		return strconv.FormatInt(int64(v), 10)
	case int64:
		return strconv.FormatInt(v, 10)
	case int32:
		return strconv.FormatInt(int64(v), 10)
	case int16:
		return strconv.FormatInt(int64(v), 10)
	case int8:
		return strconv.FormatInt(int64(v), 10)
	case uint64:
		return strconv.FormatUint(v, 10)
	case uint32:
		return strconv.FormatUint(uint64(v), 10)
	case uint16:
		return strconv.FormatUint(uint64(v), 10)
	case uint8:
		return strconv.FormatUint(uint64(v), 10)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case float32:
		return strconv.FormatFloat(float64(v), 'f', -1, 64)
	case bool:
		return strconv.FormatBool(v)
	default:
		b, err := json.Marshal(x)
		if err != nil {
			return BytesToString(b)
		}
		// 最劣的方式
		return fmt.Sprintf("%v", x)
	}
}

// 浮点数保留x位小数，最后一位四舍五入
func Round(n float64, decimals int) float64 {
	if decimals < 0 || decimals > 10 {
		panic("decimals is out of range")
	}

	if n == 0 {
		return 0
	}

	s := fmt.Sprintf("%."+strconv.Itoa(decimals)+"f", n)
	f, _ := strconv.ParseFloat(s, 64)
	return f
}

// 加密
func Sha256(s string) string {
	h := sha256.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

// 加密
func HmacSha256(s string, key []byte) string {
	h := hmac.New(sha256.New, key)
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

// base64
func Base64(src []byte) string {
	return base64.StdEncoding.EncodeToString(src)
}

// IPV4 点分十进制 转 十进制
func IPToInt(ip string) int {
	i := big.NewInt(0).SetBytes(net.ParseIP(ip).To4()).Int64()
	return int(i)
}

// IPV4 十进制 => 点分十进制
func IntToIP(ip int) string {
	return net.IPv4(byte(ip>>24), byte(ip>>16), byte(ip>>8), byte(ip)).String()
}

// 字符串长度
func StrLen(s string) int {
	return utf8.RuneCountInString(s)
}

const (
	letterDiff = 'a' - 'A'
)

// 蛇形命名，下划线格式
func SnakeCase(s string) string {
	b := make([]byte, 0, len(s)+3)
	for i := 0; i < len(s); i++ {
		v := s[i]
		if v >= 'A' && v <= 'Z' {
			if i != 0 {
				b = append(b, '_')
			}
			v += letterDiff
		}
		b = append(b, v)
	}

	return BytesToString(b)
}

// 转为JSON，禁止转义
func JsonEncode(data interface{}) ([]byte, error) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(data)
	return buffer.Bytes(), err
}

// 是否为空值
func IsEmpty[E comparable](e E) bool {
	var zero E
	return zero == e
}

// 除法。保留x位小数。最后一位四舍五入
func Div[T typex.Integer](a, b T, decimals int) float64 {
	if b == 0 {
		return 0
	}
	return Round(float64(a)/float64(b), decimals)
}

// 转字符串切片
func ToStringSlice[E typex.Integer](arr []E) []string {
	r := make([]string, len(arr))
	for i, v := range arr {
		r[i] = strconv.FormatInt(int64(v), 10)
	}
	return r
}

// 转任意类型切片
func ToAnySlice[E any](arr []E) []any {
	r := make([]any, len(arr))
	for i, v := range arr {
		r[i] = v
	}
	return r
}

// 结构体、结构体指针、结构体切片
func GetTag(v any, tag string) []string {
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
			result = append(result, GetTag(field.Type, tag)...)
		}
	}

	return result
}
