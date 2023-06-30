package candy

import (
	"math/rand"
	"time"
)

var (
	LowerLettersCharset = []rune("abcdefghijklmnopqrstuvwxyz")
	UpperLettersCharset = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	LettersCharset      = append(LowerLettersCharset, UpperLettersCharset...)
	NumbersCharset      = []rune("0123456789")
	AlphanumericCharset = append(LettersCharset, NumbersCharset...)
	SpecialCharset      = []rune("!@#$%^&*()_+-=[]{}|;':\",./<>?")
	AllCharset          = append(AlphanumericCharset, SpecialCharset...)
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// 随机数
func RandInt(max int) int {
	return rand.Intn(max)
}

// 随机字符串，自定义字符集
func RandString(size int, charset []rune) string {
	if size <= 0 {
		panic("size must be greater than 0")
	}
	if len(charset) == 0 {
		panic("charset must not be empty")
	}

	b := make([]rune, size)
	charsetLen := len(charset)
	for i := range b {
		b[i] = charset[rand.Intn(charsetLen)]
	}
	return string(b)
}
