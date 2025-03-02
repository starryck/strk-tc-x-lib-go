package xbtype

type Comparable[T any] interface {
	Compare(value T) int
}
