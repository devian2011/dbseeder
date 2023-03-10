package helper

import (
	"crypto/sha256"
	"fmt"
	"strings"
)

func InArray[T comparable](in []T, needle T) bool {
	for _, t := range in {
		if t == needle {
			return true
		}
	}

	return false
}

func SliceHash(sl []any) string {
	b := strings.Builder{}
	for _, v := range sl {
		b.WriteString(fmt.Sprintf("%v", v))
	}

	return string(sha256.New().Sum([]byte(b.String())))
}
