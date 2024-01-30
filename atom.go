package goui

import (
	"github.com/goui-org/goui/maps"
	"github.com/goui-org/goui/sets"
)

// import { Dispatch, SetStateAction, UpdateStateAction } from './hooks.js';
// import { ComponentElem } from './elem.js';
// import { callComponentFuncAndReconcile } from './reconcile.js';

// export type AtomSelector<T, R> = (state: T) => R

type selectorRecord[T comparable] struct {
	selected any
	selector func(T) any
}

type Atom[T comparable] struct {
	value     T
	subs      *sets.Set[*Node]
	selectors *maps.Map[*Node, []*selectorRecord[T]]
}

func CreateAtom[T comparable](initialValue T) *Atom[T] {
	return &Atom[T]{
		value:     initialValue,
		subs:      sets.New[*Node](),
		selectors: maps.New[*Node, []*selectorRecord[T]](),
	}
}

func (a *Atom[T]) subscribe(n *Node) {
	a.subs.Add(n)
}

func (a *Atom[T]) unsubscribe(n *Node) {
	a.subs.Delete(n)
}

func (a *Atom[T]) update(action func(T) T) {
	newVal := action(a.value)
	if newVal != a.value {
		a.value = newVal
		go a.reconcile()
	}
}

func (a *Atom[T]) reconcile() {
	for _, node := range a.subs.Slice() {
		if a.subs.Has(node) {
			callComponentFuncAndReconcile(node, node)
		}
	}
	for _, entry := range a.selectors.Slice() {
		for _, record := range entry.Value {
			if record.selected != record.selector(a.value) {
				callComponentFuncAndReconcile(entry.Key, entry.Key)
			}
		}
	}
}

// export interface Atom<T> {
//     s: T // state
//     u: Dispatch<SetStateAction<T>> // update atom state
//     r: () => void // reconcile all subscribers
//     c: Set<ComponentElem> // component subscribers
//     a: Set<Atom<any> | ReadonlyAtom<any>> // atoms subscribed to this atom
//     f: Map<ComponentElem, [any, AtomSelector<any, any>][]> // selectors subscribed to this atom
// }

// interface Deriver<T> {
//     d: AtomDerivation<T> // state derivation
// }

// export type AtomGetter = <A>(atom: Atom<A> | ReadonlyAtom<A>) => A;
// export type ReadonlyAtom<T> = Omit<Atom<T>, 'u'> & Deriver<T>;
// export type AtomDerivation<T> = (get: AtomGetter) => T;

// type CreateAtom = {
//     <T>(derivation: AtomDerivation<T>, options?: AtomOptions<T>): ReadonlyAtom<T>
//     <T>(initialValue: T, options?: AtomOptions<T>): Atom<T>
// };

// export interface AtomOptions<T> {
//     watch?: (prevState: T, newState: T) => void
// }

// export let createAtom: CreateAtom = (config: any, options?: any): any => typeof config === 'function' ? createDerivedAtom(config, options) : createStandardAtom(config, options);

// let createStandardAtom = <T>(initialValue: T, options?: AtomOptions<T>): Atom<T> => {
//     let atom: Atom<T> = {
//         s: initialValue,
//         u: action => {
//             let oldState = atom.s;
//             atom.s = typeof action === 'function' ? (action as UpdateStateAction<T>)(oldState) : action;
//             if (oldState !== atom.s) {
//                 options?.watch?.(oldState, atom.s);
//                 queueMicrotask(atom.r);
//             }
//         },
//         r: () => updateAtomSubscribers(atom),
//         c: new Set(),
//         a: new Set(),
//         f: new Map(),
//     };
//     return atom;
// };

// let createDerivedAtom = <T>(derivation: AtomDerivation<T>, options?: AtomOptions<T>): ReadonlyAtom<T> => {
//     let atom: ReadonlyAtom<T> = {
//         s: null as T,
//         d: derivation,
//         r: null as unknown as () => void,
//         c: new Set(),
//         a: new Set(),
//         f: new Map(),
//     };
//     let getter: AtomGetter = <A>(a: Atom<A> | ReadonlyAtom<A>): A => {
//         a.a.add(atom);
//         return a.s;
//     }
//     atom.s = derivation(getter);
//     atom.r = () => {
//         let oldState = atom.s;
//         atom.s = atom.d(getter);
//         if (atom.s !== oldState) {
//             options?.watch?.(oldState, atom.s);
//             updateAtomSubscribers(atom);
//         }
//     }
//     return atom;
// };

// let updateAtomSubscribers = <T>(atom: Atom<T> | ReadonlyAtom<T>): void => {
//     for (let component of [...atom.c.keys()]) {
//         if (atom.c.has(component)) {
//             callComponentFuncAndReconcile(component, component);
//         }
//     }
//     atom.a.forEach(a => a.r());
//     for (let [component, selects] of [...atom.f.entries()]) {
//         for (let i = selects.length - 1; i >= 0; i--) {
//             let [selected, selector] = selects[i];
//             if (selected !== selector(atom.s)) {
//                 callComponentFuncAndReconcile(component, component);
//             }
//         }
//     }
// };
