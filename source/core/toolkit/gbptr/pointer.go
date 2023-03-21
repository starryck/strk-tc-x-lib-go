package gbptr

func Zero[T any]() T {
	var value T
	return value
}

func Refer[T any](value T) *T {
	return &value
}

func Deref[T any](value *T) T {
	if value != nil {
		return *value
	}
	return Zero[T]()
}
