package xbctnr

func NewSet[T comparable](values ...T) Set[T] {
	set := make(Set[T], len(values))
	for _, value := range values {
		set.Add(value)
	}
	return set
}

type Set[T comparable] map[T]struct{}

func (set Set[T]) Has(value T) bool {
	_, ok := set[value]
	return ok
}

func (set Set[T]) Add(value T) {
	set[value] = struct{}{}
}

func (set Set[T]) Remove(value T) {
	delete(set, value)
}

func (set Set[T]) Clear() {
	for key := range set {
		delete(set, key)
	}
}
