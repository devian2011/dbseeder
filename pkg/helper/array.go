package helper

import (
	gbe "github.com/devian2011/go_basic_extension"
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
	return gbe.HashSlice(sl)
}
