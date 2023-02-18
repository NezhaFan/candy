package util

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"math/big"
	"net"
	"reflect"
	"strconv"
	"unsafe"
)

func BytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func StringToBytes(s string) (b []byte) {
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh.Data = sh.Data
	bh.Cap = sh.Len
	bh.Len = sh.Len
	return b
}

func FloatToString(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
}

func ToJson(v any) string {
	if v == nil {
		return ""
	}
	b, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return *(*string)(unsafe.Pointer(&b))
}

func Sha256(s string) string {
	h := sha256.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

func HmacSha256(s string, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

func Base64(src []byte) string {
	return base64.StdEncoding.EncodeToString(src)
}

// IPV4 点分十进制 转 十进制
func IpToInt(ip string) int {
	i := big.NewInt(0).SetBytes(net.ParseIP(ip).To4()).Int64()
	return int(i)
}

// IPV4 十进制 => 点分十进制
func IntToIp(ip int) string {
	return net.IPv4(byte(ip>>24), byte(ip>>16), byte(ip>>8), byte(ip)).String()
}
