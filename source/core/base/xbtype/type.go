package xbtype

func NewSet[T comparable](values ...T) map[T]struct{} {
	set := make(map[T]struct{}, len(values))
	for _, value := range values {
		set[value] = struct{}{}
	}
	return set
}

type Comparable[T any] interface {
	Compare(value T) int
}
