package goui

import (
	"reflect"
	"runtime"

	"github.com/twharmon/godom"
)

type EffectTeardown func()
type SetStateFunc[T any] func(T) T
type StateDispatcher[T any] func(SetStateFunc[T])

type effectRecord struct {
	deps []any
	td   EffectTeardown
}

type memoRecord struct {
	deps []any
	val  any
}

func usePC() uintptr {
	pc, _, _, _ := runtime.Caller(2)
	return pc
}

func UseState[T any](initialValue T) (T, StateDispatcher[T]) {
	pc := usePC()
	node := useCurrentComponent()
	states := node.getStates()
	fn := func(fn SetStateFunc[T]) {
		oldVal, ok := states.Get(pc).(T)
		if node.tornDown || !ok {
			godom.Console.Error("[GOUI] Error: unable to set state for %s after component is unmounted", node.name)
			return
		}
		newVal := fn(oldVal)
		if reflect.DeepEqual(oldVal, newVal) {
			return
		}
		states.Set(pc, newVal)
		select {
		case node.updateCh <- struct{}{}:
		default:
		}
	}
	if v := states.Get(pc); v != nil {
		return v.(T), fn
	}
	states.Set(pc, initialValue)
	return initialValue, fn
}

func UseEffect(effect func() EffectTeardown, deps ...any) {
	pc := usePC()
	node := useCurrentComponent()
	effects := node.getEffects()
	record := effects.Get(pc)
	node.pendingEffects = append(node.pendingEffects, func() {
		if record != nil {
			if reflect.DeepEqual(record.deps, deps) {
				return
			}
			record.teardown()
		}
		effects.Set(pc, &effectRecord{
			deps: deps,
			td:   effect(),
		})
	})
}

type Callback[Func any] struct {
	Invoke Func
}

func UseCallback[Func any](handlerFunc Func, deps ...any) *Callback[Func] {
	return UseMemo(func() *Callback[Func] {
		return &Callback[Func]{Invoke: handlerFunc}
	}, deps...)
}

type Ref[T any] struct {
	Current T
}

func UseRef[T any](initialValue T) *Ref[T] {
	return UseMemo(func() *Ref[T] { return &Ref[T]{Current: initialValue} })
}

func UseDeferredEffect(effect func() EffectTeardown, deps ...any) {
	first := UseRef(true)
	UseEffect(func() EffectTeardown {
		if first.Current {
			first.Current = false
			return nil
		}
		return effect()
	}, deps...)
}

func UseMemo[T any](create func() T, deps ...any) T {
	pc := usePC()
	node := useCurrentComponent()
	memos := node.getMemos()
	if record := memos.Get(pc); record != nil && reflect.DeepEqual(record.deps, deps) {
		return record.val.(T)
	}
	val := create()
	memos.Set(pc, &memoRecord{
		deps: deps,
		val:  val,
	})
	return val
}

func (r *effectRecord) teardown() {
	if r.td != nil {
		r.td()
	}
}
