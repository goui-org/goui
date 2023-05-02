package goui

import "runtime"

var components = newStore[uintptr, *Node]()

func getCallerFuncId(skip int) (uintptr, uintptr) {
	pc, _, _, _ := runtime.Caller(skip)
	return pc, runtime.FuncForPC(pc).Entry()
}

func getCallerComponent() (uintptr, *Node) {
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
	id, component := getCallerComponent()
	fn := func(fn func(T) T) {
		oldVal := component.state.get(id).(T)
		newVal := fn(oldVal)
		component.state.set(id, newVal)
		old := component.vdom
		component.vdom = component.fn(component.props)
		reconcile(old, component.vdom)
	}
	if v := component.state.get(id); v != nil {
		return v.(T), fn
	}
	component.state.set(id, initialValue)
	return initialValue, fn
}

func UseEffect(effect func() EffectTeardown, deps []any) {
	id, component := getCallerComponent()
	go func() {
		if record := component.effects.get(id); record != nil {
			if areDepsSame(record.deps, deps) {
				return
			}
			record.teardown()
		}
		component.effects.set(id, &effectRecord{
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
