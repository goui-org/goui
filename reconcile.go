package goui

import (
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
		oldElem.teardown()
		// TODO: namespace?
		getDom(oldElem).Call("replaceWith", createDom(newElem, ""))
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
	if oldElem.props != newElem.props {
		oldElem.dom.Set("data", newElem.props)
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
		var t T
		if newAttr == t {
			elem.Call("removeAttr", name)
		} else {
			elem.Set(name, newAttr)
		}
	}
}

func reconcileAttributes(oldElem *Elem, newElem *Elem) {
	oldAttrs := oldElem.props.(Attributes)
	newAttrs := newElem.props.(Attributes)

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
			newAttrs.OnClick.Invoke(newMouseEvent(args[0]))
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

func moveBefore(parent js.Value, newChdNextKey string, oldChdKey string, currDomNode js.Value, movingDomNode js.Value) {
	oldPos := movingDomNode.Get("nextSibling")
	parent.Call("insertBefore", movingDomNode, currDomNode)
	if !currDomNode.Equal(parent.Get("lastChild")) && newChdNextKey != oldChdKey {
		parent.Call("insertBefore", currDomNode, oldPos)
	}
}

func reconcileChildren(oldElem *Elem, newElem *Elem) {
	newChn := newElem.props.(Attributes).Children
	oldChn := oldElem.props.(Attributes).Children
	newLength := len(newChn)
	oldLength := len(oldChn)
	if newLength == 0 && oldLength > 0 {
		newElem.dom.Set("innerHTML", "")
		for _, ch := range oldChn {
			ch.teardown()
		}
		return
	}
	start := 0

	// prefix
	for start < newLength &&
		start < oldLength &&
		(newChn[start].key == "" || newChn[start].key == oldChn[start].key) {
		reconcile(oldChn[start], newChn[start])
		start++
	}
	if start >= newLength {
		for i := start; start < oldLength; i++ {
			oldChn[i].teardown()
			getDom(oldChn[i]).Call("remove")
		}
		return
	}

	// suffix
	oldLength--
	newLength--
	for newLength > start &&
		oldLength >= start &&
		(newChn[newLength].key == "" || newChn[newLength].key == oldChn[oldLength].key) {
		reconcile(oldChn[oldLength], newChn[newLength])
		oldLength--
		newLength--
	}

	oldMap := make(map[string]*Elem)
	for i := start; i <= oldLength; i++ {
		oldChd := oldChn[i]
		oldKey := oldChd.key
		if i >= len(newChn) {
			continue
		}
		if oldKey != "" && oldKey != newChn[i].key {
			oldMap[oldKey] = oldChd
		}
	}

	chNodes := newElem.dom.Get("childNodes")
	for start <= newLength {
		newChd := newChn[start]
		if len(oldChn) <= start {
			newElem.dom.Call("appendChild", createDom(newChd, ""))
			start++
			continue
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
		nextNewKey := ""
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
