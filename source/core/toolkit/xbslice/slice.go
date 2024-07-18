package xbslice

import (
	"github.com/starryck/x-lib-go/source/core/toolkit/xbvalue"
)

func First[T any](elems []T) T {
	if len(elems) == 0 {
		return xbvalue.Zero[T]()
	}
	return elems[0]
}

func Last[T any](elems []T) T {
	if len(elems) == 0 {
		return xbvalue.Zero[T]()
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
