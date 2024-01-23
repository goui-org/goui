package goui

type EffectTeardown func()
type Deps []any

type Callback[Func any] struct {
	invoke Func
}

type effectRecord struct {
	deps     Deps
	teardown EffectTeardown
}

func useHooks() (int, *Node) {
	node := currentElem
	cursor := node.hooksCursor
	node.hooksCursor++
	return cursor, node
}

func UseState[T comparable](initialValue T) (T, func(func(T) T)) {
	cursor, node := useHooks()
	if len(node.hooks) <= cursor {
		node.hooks = append(node.hooks, initialValue)
	}
	setState := func(update func(T) T) {
		if node.unmounted {
			panic("bad set state")
		}
		oldVal := node.hooks[cursor].(T)
		newVal := update(oldVal)
		if newVal != oldVal {
			node.hooks[cursor] = newVal
			node.queue = append(node.queue, callComponentFunc(node))
			go func() {
				if len(node.queue) > 0 {
					tip := node.queue[len(node.queue)-1]
					node.queue = node.queue[:0]
					reconcile(node.virtNode, tip)
					node.virtNode = tip
				}
			}()
		}
	}
	return node.hooks[cursor].(T), setState
}

func UseEffect(effect func() EffectTeardown, deps Deps) {
	cursor, node := useHooks()
	if len(node.hooks) <= cursor {
		record := &effectRecord{deps: deps}
		node.hooks = append(node.hooks, record)
		go func() {
			if !node.unmounted {
				record.teardown = effect()
			}
		}()
		return
	}
	record := node.hooks[cursor].(*effectRecord)
	if !areDepsEqual(deps, record.deps) {
		record.deps = deps
		go func() {
			if record.teardown != nil {
				record.teardown()
			}
			if !node.unmounted {
				record.teardown = effect()
			}
		}()
	}
}

func UseImmediateEffect(effect func() EffectTeardown, deps Deps) {
	cursor, node := useHooks()
	if len(node.hooks) <= cursor {
		node.hooks = append(node.hooks, &effectRecord{
			deps:     deps,
			teardown: effect(),
		})
		return
	}
	record := node.hooks[cursor].(*effectRecord)
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
	cursor, node := useHooks()
	if len(node.hooks) <= cursor {
		m := &memoRecord[T]{
			val:  create(),
			deps: deps,
		}
		node.hooks = append(node.hooks, m)
		return m.val
	}
	memo := node.hooks[cursor].(*memoRecord[T])
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

// func queueTask(task func()) {
// 	time.AfterFunc(0, task)
// }
