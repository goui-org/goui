package goui

import (
	"bytes"
	"reflect"
	"runtime"
)

type Deps []any

func usePC() uintptr {
	pc, _, _, _ := runtime.Caller(2)
	return pc
}

func UseState[T any](initialValue T) (T, SetStateFunc[T]) {
	pc := usePC()
	node := getCurrentNode()
	states := node.getStates()
	fn := func(fn func(T) T) {
		go func() {
			oldVal := states.get(pc).(T)
			newVal := fn(oldVal)
			if deepEqual(oldVal, newVal) {
				return
			}
			states.set(pc, newVal)
			old := node.vdom
			node.vdom = node.fn(node.props)
			reconcile(old, node.vdom)
		}()
	}
	if v := states.get(pc); v != nil {
		return v.(T), fn
	}
	states.set(pc, initialValue)
	return initialValue, fn
}

func UseEffect(effect func() EffectTeardown, deps Deps) {
	pc := usePC()
	node := getCurrentNode()
	effects := node.getEffects()
	go func() {
		if record := effects.get(pc); record != nil {
			if deepEqual(record.deps, deps) {
				return
			}
			record.teardown()
		}
		effects.set(pc, &effectRecord{
			deps: deps,
			td:   effect(),
		})
	}()
}

func UseMemo[T any](create func() T, deps Deps) T {
	pc := usePC()
	node := getCurrentNode()
	memos := node.getMemos()
	if record := memos.get(pc); record != nil && deepEqual(record.deps, deps) {
		return record.val.(T)
	}
	val := create()
	memos.set(pc, &memoRecord{
		deps: deps,
		val:  val,
	})
	return val
}

type EffectTeardown func()
type SetStateFunc[T any] func(func(T) T)

type effectRecord struct {
	deps Deps
	td   EffectTeardown
}

type memoRecord struct {
	deps Deps
	val  any
}

func (r *effectRecord) teardown() {
	if r.td != nil {
		r.td()
	}
}

func deepEqual(a any, b any) bool {
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
