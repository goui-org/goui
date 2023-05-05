package equalityutil

import (
	"bytes"
	"reflect"
)

func DeepEqual(a any, b any) bool {
	aVal := reflect.ValueOf(a)
	bVal := reflect.ValueOf(b)
	if aVal.Kind() != bVal.Kind() {
		return false
	}
	return deepValueEqual(aVal, bVal)
}

func arraySliceValueEqual(a reflect.Value, b reflect.Value) bool {
	for i := 0; i < a.Len(); i++ {
		if !deepValueEqual(a.Index(i), b.Index(i)) {
			return false
		}
	}
	return true
}

func deepValueEqual(a reflect.Value, b reflect.Value) bool {
	if a.Comparable() {
		return a.Equal(b)
	}
	switch a.Kind() {
	case reflect.Func:
		return a.UnsafePointer() == b.UnsafePointer()
	case reflect.Array:
		return arraySliceValueEqual(a, b)
	case reflect.Slice:
		if a.IsNil() != b.IsNil() {
			return false
		}
		if a.Len() != b.Len() {
			return false
		}
		if a.UnsafePointer() == b.UnsafePointer() {
			return true
		}
		if a.Type().Elem().Kind() == reflect.Uint8 {
			return bytes.Equal(a.Bytes(), b.Bytes())
		}
		return arraySliceValueEqual(a, b)
	case reflect.Struct:
		for i, n := 0, a.NumField(); i < n; i++ {
			if !deepValueEqual(a.Field(i), b.Field(i)) {
				return false
			}
		}
		return true
	case reflect.Interface:
		if a.IsNil() || b.IsNil() {
			return a.IsNil() == b.IsNil()
		}
		return deepValueEqual(a.Elem(), b.Elem())
	case reflect.Pointer:
		if a.UnsafePointer() == b.UnsafePointer() {
			return true
		}
		return deepValueEqual(a.Elem(), b.Elem())
	case reflect.Map:
		if a.IsNil() != b.IsNil() {
			return false
		}
		if a.Len() != b.Len() {
			return false
		}
		if a.UnsafePointer() == b.UnsafePointer() {
			return true
		}
		for _, k := range a.MapKeys() {
			aa := a.MapIndex(k)
			bb := b.MapIndex(k)
			if !aa.IsValid() || !bb.IsValid() || !deepValueEqual(aa, bb) {
				return false
			}
		}
		return true
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return a.Int() == b.Int()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return a.Uint() == b.Uint()
	case reflect.String:
		return a.String() == b.String()
	case reflect.Bool:
		return a.Bool() == b.Bool()
	case reflect.Float32, reflect.Float64:
		return a.Float() == b.Float()
	case reflect.Complex64, reflect.Complex128:
		return a.Complex() == b.Complex()
	default:
		return a.Interface() == b.Interface()
	}
}
