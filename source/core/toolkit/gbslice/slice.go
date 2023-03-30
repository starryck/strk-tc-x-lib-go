package gbslice

import (
	"github.com/forbot161602/pbc-golang-lib/source/core/toolkit/gbvalue"
)

func First[T any](elems []T) T {
	if len(elems) == 0 {
		return gbvalue.Zero[T]()
	}
	return elems[0]
}

func Last[T any](elems []T) T {
	if len(elems) == 0 {
		return gbvalue.Zero[T]()
	}
	return elems[len(elems)-1]
}

func Contain[T comparable](elems []T, value T) bool {
	for _, elem := range elems {
		if elem == value {
			return true
		}
	}
	return false
}
