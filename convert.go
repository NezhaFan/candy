package candy

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"strconv"
)

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

// 蛇形命名，下划线格式
func ToSnakeCase(s string) string {
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

	return string(b)
}

func ToString(i any) string {
	switch x := i.(type) {
	case nil:
		return ""
	case string:
		return x
	case float64:
		return strconv.FormatFloat(x, 'f', -1, 64)
	case float32:
		return strconv.FormatFloat(float64(x), 'f', -1, 32)
	case int:
		return strconv.FormatInt(int64(x), 10)
	case int64:
		return strconv.FormatInt(x, 10)
	case int32:
		return strconv.FormatInt(int64(x), 10)
	case int16:
		return strconv.FormatInt(int64(x), 10)
	case int8:
		return strconv.FormatInt(int64(x), 10)
	case uint:
		return strconv.FormatUint(uint64(x), 10)
	case uint64:
		return strconv.FormatUint(uint64(x), 10)
	case uint32:
		return strconv.FormatUint(uint64(x), 10)
	case uint16:
		return strconv.FormatUint(uint64(x), 10)
	case uint8:
		return strconv.FormatUint(uint64(x), 10)
	case bool:
		return strconv.FormatBool(x)
	case []byte:
		return string(x)
	default:
		b, _ := json.Marshal(i)
		return string(b)
	}
}

func ToInt(i any) int {
	switch x := i.(type) {
	case nil:
		return 0
	case string:
		y, _ := strconv.Atoi(x)
		return y
	case float64:
		return int(x)
	case float32:
		return int(x)
	case int:
		return x
	case int64:
		return int(x)
	case int32:
		return int(x)
	case int16:
		return int(x)
	case int8:
		return int(x)
	case uint:
		return int(x)
	case uint64:
		return int(x)
	case uint32:
		return int(x)
	case uint16:
		return int(x)
	case uint8:
		return int(x)
	case bool:
		if x {
			return 1
		}
		return 0
	default:
		return 0
	}
}

func ToUInt(i any) uint {
	switch x := i.(type) {
	case nil:
		return 0
	case string:
		y, _ := strconv.Atoi(x)
		return uint(y)
	case float64:
		return uint(x)
	case float32:
		return uint(x)
	case int:
		return uint(x)
	case int64:
		return uint(x)
	case int32:
		return uint(x)
	case int16:
		return uint(x)
	case int8:
		return uint(x)
	case uint:
		return x
	case uint64:
		return uint(x)
	case uint32:
		return uint(x)
	case uint16:
		return uint(x)
	case uint8:
		return uint(x)
	case bool:
		if x {
			return 1
		}
		return 0
	default:
		return 0
	}
}

func ToFloat32(i any) float32 {
	switch x := i.(type) {
	case nil:
		return 0
	case float32:
		return x
	case float64:
		return float32(x)
	case int64:
		return float32(x)
	case int32:
		return float32(x)
	case int16:
		return float32(x)
	case int8:
		return float32(x)
	case uint:
		return float32(x)
	case uint64:
		return float32(x)
	case uint32:
		return float32(x)
	case uint16:
		return float32(x)
	case uint8:
		return float32(x)
	case bool:
		if x {
			return 1
		}
		return 0
	case string:
		v, _ := strconv.ParseFloat(x, 32)
		return float32(v)
	default:
		return 0
	}
}

func ToFloat64(i any) float64 {
	switch x := i.(type) {
	case nil:
		return 0
	case float32:
		return float64(x)
	case float64:
		return x
	case int64:
		return float64(x)
	case int32:
		return float64(x)
	case int16:
		return float64(x)
	case int8:
		return float64(x)
	case uint:
		return float64(x)
	case uint64:
		return float64(x)
	case uint32:
		return float64(x)
	case uint16:
		return float64(x)
	case uint8:
		return float64(x)
	case bool:
		if x {
			return 1
		}
		return 0
	case string:
		v, _ := strconv.ParseFloat(x, 32)
		return float64(v)
	default:
		return 0
	}
}
