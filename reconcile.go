package goui

import (
	"reflect"
	"syscall/js"
)

func getDom(elem *Elem) js.Value {
	for elem.virt != nil {
		elem = elem.virt
	}
	return elem.dom
}

func reconcile(oldElem *Elem, newElem *Elem) {
	newElem.unmounted = false
	if oldElem.tag != newElem.tag || oldElem.ptr != newElem.ptr {
		oldDom := getDom(oldElem)
		oldElem.teardown()
		// TODO: namespace?
		oldDom.Call("replaceWith", createDom(newElem, ""))
		return
	}
	if oldElem.render != nil {
		reconcileComponents(oldElem, newElem)
		return
	}
	newElem.dom = oldElem.dom
	if oldElem.tag != "" {
		reconcileVdomElems(oldElem, newElem)
	} else {
		reconcileTextElems(oldElem, newElem)
	}
}

func reconcileVdomElems(oldElem *Elem, newElem *Elem) {
	reconcileAttributes(oldElem, newElem)
	reconcileReference(oldElem, newElem)
	reconcileChildren(oldElem, newElem)
}

func reconcileTextElems(oldElem *Elem, newElem *Elem) {
	if oldElem.attrs != newElem.attrs {
		oldElem.dom.Set("data", newElem.attrs)
	}
}

func reconcileComponents(oldElem *Elem, newElem *Elem) {
	newElem.hooks = oldElem.hooks
	if oldElem.memo != nil && newElem.memo != nil && areDepsEqual(oldElem.memo, newElem.memo) {
		newElem.virt = oldElem.virt
		return
	}
	oldElem.queue = nil
	callComponentFuncAndReconcile(oldElem, newElem)
}

func reconcileAttribute[T comparable](oldAttr T, newAttr T, name string, elem js.Value) {
	if oldAttr != newAttr {
		if reflect.ValueOf(newAttr).IsZero() {
			elem.Call("removeAttr", name)
		} else {
			elem.Set(name, newAttr)
		}
	}
}

func reconcileAttributes(oldElem *Elem, newElem *Elem) {
	oldAttrs := oldElem.attrs.(*Attributes)
	newAttrs := newElem.attrs.(*Attributes)

	if oldAttrs.Class != newAttrs.Class {
		if newAttrs.Class == "" {
			oldElem.dom.Call("removeAttribute", "class")
		} else {
			oldElem.dom.Set("className", newAttrs.Class)
		}
	}
	if oldAttrs.AriaHidden != newAttrs.AriaHidden {
		if !newAttrs.AriaHidden {
			oldElem.dom.Call("removeAttribute", "aria-hidden")
		} else {
			oldElem.dom.Set("ariaHidden", newAttrs.AriaHidden)
		}
	}
	reconcileAttribute(oldAttrs.Style, newAttrs.Style, "style", oldElem.dom)
	reconcileAttribute(oldAttrs.ID, newAttrs.ID, "id", oldElem.dom)
	reconcileAttribute(oldAttrs.Disabled, newAttrs.Disabled, "disabled", oldElem.dom)
	reconcileAttribute(oldAttrs.Value, newAttrs.Value, "value", oldElem.dom)
	if oldAttrs.OnClick != newAttrs.OnClick {
		oldElem.dom.Set("onclick", js.FuncOf(func(_ js.Value, args []js.Value) any {
			newAttrs.OnClick.invoke(newMouseEvent(args[0]))
			return nil
		}))
	}
}

func reconcileReference(oldElem *Elem, newElem *Elem) {
	if newElem.ref != nil {
		newElem.ref.Value = oldElem.dom
	} else if oldElem.ref != nil {
		oldElem.ref.Value = js.Null()
	}
}

func callComponentFuncAndReconcile(oldElem *Elem, newElem *Elem) {
	newElemVdom := callComponentFunc(newElem)
	reconcile(oldElem.virt, newElemVdom)
	newElem.virt = newElemVdom
}

func moveBefore(parent js.Value, newChdNextKey any, oldChdKey any, currDomNode js.Value, movingDomNode js.Value) {
	oldPos := movingDomNode.Get("nextSibling")
	parent.Call("insertBefore", movingDomNode, currDomNode)
	if newChdNextKey != oldChdKey && !currDomNode.Equal(parent.Get("lastChild")) {
		parent.Call("insertBefore", currDomNode, oldPos)
	}
}

func reconcileChildren(oldElem *Elem, newElem *Elem) {
	newChn := newElem.attrs.(*Attributes).Children
	oldChn := oldElem.attrs.(*Attributes).Children
	newLength := len(newChn)
	oldLength := len(oldChn)
	if newLength == 0 && oldLength > 0 {
		// newElem.dom.Set("innerHTML", "")
		for _, ch := range oldChn {
			getDom(ch).Call("remove")
			ch.teardown()
		}
		return
	}
	start := 0

	// prefix
	for start < newLength && start < oldLength {
		o := oldChn[start]
		n := newChn[start]
		if n.key == nil || n.key == o.key {
			reconcile(o, n)
		} else {
			break
		}
		start++
	}
	if start >= newLength {
		for i := start; start < oldLength; i++ {
			getDom(oldChn[i]).Call("remove")
			oldChn[i].teardown()
		}
		return
	}

	// suffix
	oldLength--
	newLength--
	for newLength > start && oldLength >= start {
		o := oldChn[oldLength]
		n := newChn[newLength]
		if n.key == nil || n.key == o.key {
			reconcile(o, n)
		} else {
			break
		}
		oldLength--
		newLength--
	}

	oldMap := make(map[any]*Elem)
	for i := start; i <= oldLength; i++ {
		oldChd := oldChn[i]
		oldKey := oldChd.key
		noMoreNewChn := false
		if i >= len(newChn) {
			noMoreNewChn = true
		}
		if oldKey != nil && (noMoreNewChn || oldKey != newChn[i].key) {
			oldMap[oldKey] = oldChd
		}
	}

	chNodes := newElem.dom.Get("childNodes")
	for start <= newLength {
		newChd := newChn[start]
		if len(oldChn) <= start {
			doms := make([]any, newLength-start+1)
			for i := start; i <= newLength; i++ {
				doms[i-start] = createDom(newChn[i], "")
			}
			newElem.dom.Call("append", doms...)
			// start++
			break
		}
		oldChd := oldChn[start]
		newKey := newChd.key
		if oldChd.key == newKey {
			reconcile(oldChd, newChd)
			start++
			continue
		}
		mappedOld := oldMap[newKey]
		chdDom := chNodes.Index(start)
		var nextNewKey any
		if len(newChn) > start+1 {
			nextNewKey = newChn[start+1].key
		}
		if mappedOld != nil {
			oldDom := getDom(mappedOld)
			if !chdDom.Equal(oldDom) {
				moveBefore(newElem.dom, nextNewKey, oldChd.key, chdDom, oldDom)
			}
			reconcile(mappedOld, newChd)
			delete(oldMap, newKey)
		} else {
			moveBefore(newElem.dom, nextNewKey, oldChd.key, chdDom, createDom(newChd, ""))
		}
		start++
	}

	for _, elem := range oldMap {
		getDom(elem).Call("remove")
		elem.teardown()
	}
}
