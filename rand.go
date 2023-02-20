package util

import (
	"math/rand"
	"time"
)

const (
	Number      = "0123456789"
	LowerLetter = "abcdefghijklmnopqrstuvwxyz"
	UpperLetter = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// 随机数
func RandInt(max int) int {
	return rand.Intn(max)
}

// 随机字符串，自定义字符范围
func RandString(length int, letter string) string {
	b := make([]byte, length)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	n := len(letter)
	for i := 0; i < length; i++ {
		b[i] = letter[r.Intn(n)]
	}
	return string(b)
}
