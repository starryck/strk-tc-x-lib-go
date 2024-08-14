package xbtype

type Comparable[T any] interface {
	Compare(comparable T) int
}
