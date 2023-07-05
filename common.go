package candy

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"math/big"
	"net"
	"reflect"
	"strconv"
	"unicode/utf8"
	"unsafe"
)

// bytes转字符串
func BytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// 字符串转bytes
func StringToBytes(s string) (b []byte) {
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh.Data = sh.Data
	bh.Cap = sh.Len
	bh.Len = sh.Len
	return b
}

// 浮点数转字符串
func FloatToString(n float64) string {
	return strconv.FormatFloat(float64(n), 'f', -1, 64)
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
func RuneLength(s string) int {
	return utf8.RuneCountInString(s)
}

// 是否为空值
func IsEmpty[E comparable](e E) bool {
	var zero E
	return zero == e
}

// 最大值
func Max[T Integer](n1 T, ns ...T) T {
	max := n1
	for _, n := range ns {
		if max < n {
			max = n
		}
	}

	return max
}

// 最小值
func Min[T Integer](n1 T, ns ...T) T {
	min := n1
	for _, n := range ns {
		if min > n {
			min = n
		}
	}

	return min
}

// 除法。保留x位小数。最后一位四舍五入
func Div[T Integer](a, b T, decimals int) float64 {
	return Round(float64(a)/float64(b), decimals)
}
