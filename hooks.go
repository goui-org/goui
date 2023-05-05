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
		oldVal := states.get(pc).(T)
		newVal := fn(oldVal)
		if equalityutil.DeepEqual(oldVal, newVal) {
			return
		}
		states.set(pc, newVal)
		old := node.vdom
		node.vdom = node.fn(node.props)
		reconcile(old, node.vdom)
	}
	if v := states.get(pc); v != nil {
		return v.(T), fn
	}
	states.set(pc, initialValue)
	return initialValue, fn
}

func UseEffect(effect func() EffectTeardown, deps Deps) {
	pc := usePC()
	node := useCurrentComponent()
	effects := node.getEffects()
	go func() {
		if record := effects.get(pc); record != nil {
			if equalityutil.DeepEqual(record.deps, deps) {
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
	node := useCurrentComponent()
	memos := node.getMemos()
	if record := memos.get(pc); record != nil && equalityutil.DeepEqual(record.deps, deps) {
		return record.val.(T)
	}
	val := create()
	memos.set(pc, &memoRecord{
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
