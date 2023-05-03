package goui

import (
	"fmt"
	"runtime"
)

var components = newStore[uintptr, *Node]()

func getCallerFuncId(skip int) (uintptr, uintptr) {
	pc, _, _, _ := runtime.Caller(skip)
	return pc, runtime.FuncForPC(pc).Entry()
}

func usePCAndComponentID() (uintptr, *Node) {
	i := 1
	for {
		pc, name := getCallerFuncId(i)
		if component := components.get(name); component != nil {
			return pc, component
		}
		i++
		if i > 1000 {
			panic("component not found")
		}
	}
}

func UseState[T any](initialValue T) (T, func(func(T) T)) {
	pc, component := usePCAndComponentID()
	fn := func(fn func(T) T) {
		go func() {
			oldVal := component.state.get(pc).(T)
			newVal := fn(oldVal)
			if fmt.Sprintf("%v", oldVal) == fmt.Sprintf("%v", newVal) {
				return
			}
			component.state.set(pc, newVal)
			old := component.vdom
			component.vdom = component.fn(component.props)
			reconcile(old, component.vdom)
		}()
	}
	if v := component.state.get(pc); v != nil {
		return v.(T), fn
	}
	component.state.set(pc, initialValue)
	return initialValue, fn
}

func UseEffect(effect func() EffectTeardown, deps []any) {
	pc, component := usePCAndComponentID()
	go func() {
		if record := component.effects.get(pc); record != nil {
			if areDepsSame(record.deps, deps) {
				return
			}
			record.teardown()
		}
		component.effects.set(pc, &effectRecord{
			deps: deps,
			td:   effect(),
		})
	}()
}

type EffectTeardown func()

type effectRecord struct {
	deps []any
	td   EffectTeardown
}

func (r *effectRecord) teardown() {
	if r.td != nil {
		r.td()
	}
}

func areDepsSame(a []any, b []any) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
