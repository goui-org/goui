package goui

import (
	"time"
)

type EffectTeardown func()
type Deps []any

type effectRecord struct {
	deps     Deps
	teardown EffectTeardown
}

func useHooks() (int, *Elem) {
	elem := currentElem
	cursor := elem.hooksCursor
	elem.hooksCursor++
	return cursor, elem
}

func UseState[T comparable](initialValue T) (T, func(func(T) T)) {
	cursor, elem := useHooks()
	if len(elem.hooks) <= cursor {
		elem.hooks = append(elem.hooks, initialValue)
	}
	setState := func(update func(T) T) {
		if elem.unmounted {
			panic("bad set state")
		}
		oldVal := elem.hooks[cursor].(T)
		newVal := update(oldVal)
		if newVal != oldVal {
			elem.hooks[cursor] = newVal
			elem.queue = append(elem.queue, callComponentFunc(elem))
			queueTask(func() {
				if len(elem.queue) > 0 {
					tip := elem.queue[len(elem.queue)-1]
					elem.queue = nil
					reconcile(elem.virt, tip)
					elem.virt = tip
				}
			})
		}
	}
	return elem.hooks[cursor].(T), setState
}

func UseEffect(effect func() EffectTeardown, deps Deps) {
	cursor, elem := useHooks()
	if len(elem.hooks) <= cursor {
		record := &effectRecord{deps: deps}
		elem.hooks = append(elem.hooks, record)
		queueTask(func() {
			if !elem.unmounted {
				record.teardown = effect()
			}
		})
		return
	}
	record := elem.hooks[cursor].(*effectRecord)
	if !areDepsEqual(deps, record.deps) {
		record.deps = deps
		queueTask(func() {
			if record.teardown != nil {
				record.teardown()
			}
			if !elem.unmounted {
				record.teardown = effect()
			}
		})
	}
}

func UseImmediateEffect(effect func() EffectTeardown, deps Deps) {
	cursor, elem := useHooks()
	if len(elem.hooks) <= cursor {
		elem.hooks = append(elem.hooks, &effectRecord{
			deps:     deps,
			teardown: effect(),
		})
		return
	}
	record := elem.hooks[cursor].(*effectRecord)
	if !areDepsEqual(deps, record.deps) {
		if record.teardown != nil {
			record.teardown()
		}
		record.deps = deps
		record.teardown = effect()
	}
}

type memoRecord[T any] struct {
	deps Deps
	val  T
}

func UseMemo[T any](create func() T, deps Deps) T {
	cursor, elem := useHooks()
	if len(elem.hooks) <= cursor {
		m := &memoRecord[T]{
			val:  create(),
			deps: deps,
		}
		elem.hooks = append(elem.hooks, m)
		return m.val
	}
	memo := elem.hooks[cursor].(*memoRecord[T])
	if !areDepsEqual(deps, memo.deps) {
		memo.deps = deps
		memo.val = create()
	}
	return memo.val
}

func UseCallback[Func any](handlerFunc Func, deps Deps) *Callback[Func] {
	return UseMemo(func() *Callback[Func] {
		return &Callback[Func]{invoke: handlerFunc}
	}, deps)
}

type Ref[T any] struct {
	Value T
}

func UseRef[T any](initialValue T) *Ref[T] {
	return UseMemo[*Ref[T]](func() *Ref[T] { return &Ref[T]{Value: initialValue} }, Deps{})
}

// let useAtomSubscription = <T>(atom: Atom<T> | ReadonlyAtom<T>) => {
//     let elem = current.e!;
//     useImmediateEffect(() => {
//         atom.c.add(elem);
//         return () => atom.c.delete(elem);
//     }, [elem]);
// };

// export let useAtom = <T>(atom: Atom<T>): [T, Dispatch<SetStateAction<T>>] => {
//     useAtomSubscription(atom);
//     return [atom.s, atom.u];
// };

// export let useAtomSetter = <T>(atom: Atom<T>): Dispatch<SetStateAction<T>> => atom.u;

// export let useAtomValue = <T>(atom: Atom<T> | ReadonlyAtom<T>): T => {
//     useAtomSubscription(atom);
//     return atom.s;
// };

// export let useAtomSelector = <T, R>(atom: Atom<T> | ReadonlyAtom<T>, selector: (state: T) => R): R => {
//     let elem = current.e!;
//     useImmediateEffect(() => {
//         let selected = selector(atom.s);
//         let selects = atom.f.get(elem);
//         if (!selects) {
//             atom.f.set(elem, [[selected, selector]]);
//         } else {
//             selects.push([selected, selector]);
//         }
//         return () => atom.f.delete(elem);
//     }, [elem, selector]);
//     return selector(atom.s);
// };

func queueTask(task func()) {
	time.AfterFunc(0, task)
}
