package querybuilder

import (
	"reflect"
	"time"
)

func isSupportedType(t reflect.Type) bool {
	if t.Kind() == reflect.Ptr {
		return isSupportedType(t.Elem())
	}
	supportedKinds := []reflect.Kind{
		// from type.go
		// reflect.Invalid,
		reflect.Bool,
		reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64,
		reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64,
		reflect.Uintptr,
		reflect.Float32,
		reflect.Float64,
		// reflect.Complex64,
		// reflect.Complex128,
		// reflect.Array,
		// reflect.Chan,
		// reflect.Func,
		// reflect.Interface,
		// reflect.Map,
		// reflect.Ptr,
		// reflect.Slice,
		reflect.String,
		// reflect.Struct,
		// reflect.UnsafePointer,
	}
	supportedTypes := []reflect.Type{
		reflect.TypeOf(time.Time{}),
	}
	for i := range supportedKinds {
		if t.Kind() == supportedKinds[i] {
			return true
		}
	}
	for i := range supportedTypes {
		if t == supportedTypes[i] {
			return true
		}
	}

	return false
}
