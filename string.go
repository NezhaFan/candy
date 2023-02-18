package util

import (
	"strings"
)

func PadEnd(source string, size int, padStr string) string {
	len1 := len(source)
	len2 := len(padStr)

	if len1 >= size || padStr == "" {
		return source
	}

	if len1+len2 >= size {
		return source + padStr[:size-len1]
	}

	fill := strings.Repeat(padStr, (size-len1)/len2+1)[:size-len1]
	return source + fill[:size-len1]
}

func PadStart(source string, size int, padStr string) string {
	len1 := len(source)
	len2 := len(padStr)

	if len1 >= size || padStr == "" {
		return source
	}

	if len1+len2 >= size {
		return padStr[:size-len1] + source
	}

	fill := strings.Repeat(padStr, (size-len1)/len2+1)[:size-len1]
	return fill[:size-len1] + source
}
