package gbslice

import "fmt"

func Last[T any](values []T) T {
	if len(values) == 0 {
		panic(fmt.Sprintf("Values `%v` must be a nonempty slice.", values))
	}
	return values[len(values)-1]
}
