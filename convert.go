package candy

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"net"
	"strconv"
	"unsafe"
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

// base64
func Base64(src []byte) string {
	return base64.StdEncoding.EncodeToString(src)
}

// 数字转62进制字符串
func EncodeNumber(n uint64) string {
	const base = 62
	const charset = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	var result []byte
	for n > 0 {
		r := n % base
		n /= base
		result = append(result, charset[r])
	}
	return BytesToString(result)
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

// IPV4 点分十进制 转 十进制
func Ip2long(ip string) int {
	i := big.NewInt(0).SetBytes(net.ParseIP(ip).To4()).Int64()
	return int(i)
}

// IPV4 十进制 => 点分十进制
func Long2ip(ip int) string {
	return net.IPv4(byte(ip>>24), byte(ip>>16), byte(ip>>8), byte(ip)).String()
}

// 蛇形命名，下划线格式
func SnakeCase(s string) string {
	b := make([]byte, 0, len(s)+3)
	for i := 0; i < len(s); i++ {
		v := s[i]
		if v >= 'A' && v <= 'Z' {
			if i != 0 {
				b = append(b, '_')
			}
			v += 'a' - 'A'
		}
		b = append(b, v)
	}

	return BytesToString(b)
}

func AnyToString(v any) string {
	switch x := v.(type) {
	case string:
		return x
	case []byte:
		return BytesToString(x)
	case int:
		return strconv.Itoa(x)
	case int8:
		return strconv.Itoa(int(x))
	case int16:
		return strconv.Itoa(int(x))
	case int32:
		return strconv.Itoa(int(x))
	case int64:
		return strconv.Itoa(int(x))
	case uint:
		return strconv.Itoa(int(x))
	case uint8:
		return strconv.Itoa(int(x))
	case uint16:
		return strconv.Itoa(int(x))
	case uint32:
		return strconv.Itoa(int(x))
	case uint64:
		return strconv.FormatUint(x, 64)
	case float32:
		return strconv.FormatFloat(float64(x), 'f', -1, 64)
	case float64:
		return strconv.FormatFloat(x, 'f', -1, 64)
	case bool:
		return strconv.FormatBool(x)
	case error:
		return x.Error()
	default:
		b, err := json.Marshal(v)
		if err != nil {
			return fmt.Sprintf("%v", v)
		}
		return BytesToString(b)
	}
}
