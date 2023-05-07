package goui

import (
	"runtime"

	"github.com/twharmon/goui/utils/equalityutil"
)

type Deps []any
type EffectTeardown func()
type SetStateFunc[T any] func(T) T
type StateDispatcher[T any] func(SetStateFunc[T])

type effectRecord struct {
	deps Deps
	td   EffectTeardown
}

type memoRecord struct {
	deps Deps
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
		oldVal := states.Get(pc).(T)
		newVal := fn(oldVal)
		if equalityutil.DeepEqual(oldVal, newVal) {
			return
		}
		states.Set(pc, newVal)
		node.updateCh <- struct{}{}
	}
	if v := states.Get(pc); v != nil {
		return v.(T), fn
	}
	states.Set(pc, initialValue)
	return initialValue, fn
}

func UseEffect(effect func() EffectTeardown, deps Deps) {
	pc := usePC()
	node := useCurrentComponent()
	effects := node.getEffects()
	node.pendingEffects = append(node.pendingEffects, func() {
		if record := effects.Get(pc); record != nil {
			if equalityutil.DeepEqual(record.deps, deps) {
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

type Ref[T any] struct {
	Current T
}

func UseRef[T any](initialValue T) *Ref[T] {
	return UseMemo(func() *Ref[T] { return &Ref[T]{Current: initialValue} }, Deps{})
}

func UseDeferredEffect(effect func() EffectTeardown, deps Deps) {
	first := UseRef(true)
	UseEffect(func() EffectTeardown {
		if first.Current {
			first.Current = false
			return nil
		}
		return effect()
	}, deps)
}

func UseMemo[T any](create func() T, deps Deps) T {
	pc := usePC()
	node := useCurrentComponent()
	memos := node.getMemos()
	if record := memos.Get(pc); record != nil && equalityutil.DeepEqual(record.deps, deps) {
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
