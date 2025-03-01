package xbreflect

import "reflect"

type Value = reflect.Value

func IsNil(input any) bool {
	return isNil(reflect.ValueOf(input))
}

func isNil(value Value) bool {
	switch kind := value.Kind(); kind {
	case reflect.Invalid:
		return true
	case reflect.Map, reflect.Slice, reflect.Chan, reflect.Func:
		return value.IsNil()
	case reflect.Ptr, reflect.Interface:
		return isNil(value.Elem())
	}
	return false
}

func IsAnyNil(inputs ...any) bool {
	for _, input := range inputs {
		if IsNil(input) {
			return true
		}
	}
	return false
}

func AreAllNil(inputs ...any) bool {
	if len(inputs) == 0 {
		return false
	}
	for _, input := range inputs {
		if !IsNil(input) {
			return false
		}
	}
	return true
}

// Note: Unexported struct fields would be ignored.
func IsAnyMemberNil(input any) bool {
	return isAnyMemberNil(reflect.ValueOf(input))
}

func isAnyMemberNil(value Value) bool {
	switch kind := value.Kind(); kind {
	case reflect.Invalid:
		return true
	case reflect.Chan, reflect.Func:
		return value.IsNil()
	case reflect.Ptr, reflect.Interface:
		return isAnyMemberNil(value.Elem())
	case reflect.Map:
		if value.IsNil() {
			return true
		}
		if value.Len() == 0 {
			return false
		}
		iter := value.MapRange()
		for iter.Next() {
			if isNil(iter.Value()) {
				return true
			}
		}
		return false
	case reflect.Slice:
		if value.IsNil() {
			return true
		}
		fallthrough
	case reflect.Array:
		length := value.Len()
		if length == 0 {
			return false
		}
		for i := range length {
			if isNil(value.Index(i)) {
				return true
			}
		}
		return false
	case reflect.Struct:
		length := value.NumField()
		if length == 0 {
			return false
		}
		for i := range length {
			elem := value.Field(i)
			if elem.CanInterface() && isNil(elem) {
				return true
			}
		}
		return false
	}
	return false
}

// Note: Unexported struct fields would be ignored.
func AreAllMembersNil(input any) bool {
	return areAllMembersNil(reflect.ValueOf(input))
}

func areAllMembersNil(value Value) bool {
	switch kind := value.Kind(); kind {
	case reflect.Invalid:
		return true
	case reflect.Chan, reflect.Func:
		return value.IsNil()
	case reflect.Ptr, reflect.Interface:
		return areAllMembersNil(value.Elem())
	case reflect.Map:
		if value.IsNil() {
			return true
		}
		if value.Len() == 0 {
			return false
		}
		iter := value.MapRange()
		for iter.Next() {
			if !isNil(iter.Value()) {
				return false
			}
		}
		return true
	case reflect.Slice:
		if value.IsNil() {
			return true
		}
		fallthrough
	case reflect.Array:
		length := value.Len()
		if length == 0 {
			return false
		}
		for i := range length {
			if !isNil(value.Index(i)) {
				return false
			}
		}
		return true
	case reflect.Struct:
		length := value.NumField()
		if length == 0 {
			return false
		}
		lengthOfUnexportedFields := 0
		for i := range length {
			elem := value.Field(i)
			if !elem.CanInterface() {
				lengthOfUnexportedFields++
				continue
			}
			if !isNil(elem) {
				return false
			}
		}
		if length == lengthOfUnexportedFields {
			return false
		}
		return true
	}
	return false
}

func Destine(input any) any {
	return destine(reflect.ValueOf(input))
}

func destine(value Value) any {
	switch kind := value.Kind(); kind {
	case reflect.Invalid:
		return nil
	case reflect.Ptr:
		return destine(value.Elem())
	}
	return value.Interface()
}
