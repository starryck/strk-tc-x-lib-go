package xbslice

import (
	"math"

	"github.com/starryck/strk-tc-x-lib-go/source/core/toolkit/xbvalue"
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

func Copy[T any](elems []T, lower, upper, step int) []T {
	size := int(math.Ceil(float64(upper-lower) / float64(step)))
	slice := make([]T, size)
	for i := range size {
		slice[i] = elems[lower+i*step]
	}
	return slice
}

func Defaults[T any](size int, value T) []T {
	slice := make([]T, size)
	for i := range size {
		slice[i] = value
	}
	return slice
}

func ToInterfaces[T any](elems []T) []any {
	slice := make([]any, len(elems))
	for i, elem := range elems {
		slice[i] = elem
	}
	return slice
}
