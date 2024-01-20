package goui

import (
	"reflect"
	"syscall/js"
)

func getDom(node *Node) int {
	for node.virtNode != nil {
		node = node.virtNode
	}
	return node.dom
}

func reconcile(oldNode *Node, newNode *Node) {
	newNode.unmounted = false
	if oldNode.tag != newNode.tag || oldNode.ptr != newNode.ptr {
		// oldDom := getDom(oldNode)
		oldNode.teardown()
		// TODO: namespace?
		// oldDom.Call("replaceWith", createDom(newNode, ""))
		return
	}
	if oldNode.render != nil {
		reconcileComponents(oldNode, newNode)
		return
	}
	newNode.dom = oldNode.dom
	if oldNode.tag != "" {
		reconcileVdomElems(oldNode, newNode)
	} else {
		reconcileTextElems(oldNode, newNode)
	}
}

func reconcileVdomElems(oldNode *Node, newNode *Node) {
	reconcileAttributes(oldNode, newNode)
	reconcileReference(oldNode, newNode)
	reconcileChildren(oldNode, newNode)
}

func reconcileTextElems(oldNode *Node, newNode *Node) {
	if oldNode.attrs != newNode.attrs {
		// oldNode.dom.Set("data", newNode.attrs)
		set(newNode.dom, "data", newNode.attrs.(string))
	}
}

func reconcileComponents(oldNode *Node, newNode *Node) {
	newNode.hooks = oldNode.hooks
	if oldNode.memo != nil && newNode.memo != nil && areDepsEqual(oldNode.memo, newNode.memo) {
		newNode.virtNode = oldNode.virtNode
		return
	}
	oldNode.queue = nil
	callComponentFuncAndReconcile(oldNode, newNode)
}

func reconcileAttribute[T comparable](oldAttr T, newAttr T, name string, jsVal int) {
	if oldAttr != newAttr {
		if reflect.ValueOf(newAttr).IsZero() {
			// jsVal.Call("removeAttribute", name)
		} else {
			// jsVal.Set(name, newAttr)
		}
	}
}

func reconcileAttributes(oldNode *Node, newNode *Node) {
	oldAttrs := oldNode.attrs.(*Attributes)
	newAttrs := newNode.attrs.(*Attributes)

	if oldAttrs.Class != newAttrs.Class {
		if newAttrs.Class == "" {
			// oldNode.dom.Call("removeAttribute", "class")
			setClass(oldNode.dom, "")
		} else {
			// oldNode.dom.Set("className", newAttrs.Class)
			setClass(oldNode.dom, newAttrs.Class)
		}
	}
	// if oldAttrs.AriaHidden != newAttrs.AriaHidden {
	// 	if !newAttrs.AriaHidden {
	// 		oldNode.dom.Call("removeAttribute", "aria-hidden")
	// 	} else {
	// 		oldNode.dom.Set("ariaHidden", newAttrs.AriaHidden)
	// 	}
	// }
	reconcileAttribute(oldAttrs.Style, newAttrs.Style, "style", oldNode.dom)
	reconcileAttribute(oldAttrs.ID, newAttrs.ID, "id", oldNode.dom)
	reconcileAttribute(oldAttrs.Disabled, newAttrs.Disabled, "disabled", oldNode.dom)
	reconcileAttribute(oldAttrs.Value, newAttrs.Value, "value", oldNode.dom)
	if oldAttrs.OnClick != newAttrs.OnClick {
		// oldNode.dom.Set("onclick", js.FuncOf(func(_ js.Value, args []js.Value) any {
		// 	newAttrs.OnClick.invoke(newMouseEvent(args[0]))
		// 	return nil
		// }))
		oldNode.setEventListener("onclick", func(_ js.Value, args []js.Value) any {
			newAttrs.OnClick.invoke(newMouseEvent(args[0]))
			return nil
		})
	}
}

func reconcileReference(oldNode *Node, newNode *Node) {
	// if newNode.ref != nil {
	// 	newNode.ref.Value = oldNode.dom
	// } else if oldNode.ref != nil {
	// 	oldNode.ref.Value = js.Null()
	// }
}

func callComponentFuncAndReconcile(oldNode *Node, newNode *Node) {
	newElemVdom := callComponentFunc(newNode)
	reconcile(oldNode.virtNode, newElemVdom)
	newNode.virtNode = newElemVdom
}

func moveBefore(parent js.Value, newChdNextKey any, oldChdKey any, currDomNode js.Value, movingDomNode js.Value) {
	oldPos := movingDomNode.Get("nextSibling")
	parent.Call("insertBefore", movingDomNode, currDomNode)
	if newChdNextKey != oldChdKey && !currDomNode.Equal(parent.Get("lastChild")) {
		parent.Call("insertBefore", currDomNode, oldPos)
	}
}

func reconcileChildren(oldNode *Node, newNode *Node) {
	newChn := newNode.attrs.(*Attributes).Children
	oldChn := oldNode.attrs.(*Attributes).Children
	newLength := len(newChn)
	oldLength := len(oldChn)
	if newLength == 0 && oldLength > 0 {
		// newNode.dom.Set("innerHTML", nil)
		clearChildren(newNode.dom)
		for _, ch := range oldChn {
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
			// getDom(oldChn[i]).Call("remove")
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

	oldMap := make(map[any]*Node)
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

	// chNodes := newNode.dom.Get("childNodes")
	// chNodes := global.Get("elements").Index(newNode.dom).Get(chNodes)
	for start <= newLength {
		newChd := newChn[start]
		if len(oldChn) <= start {
			// doms := make([]any, newLength-start+1)
			for i := start; i <= newLength; i++ {
				// doms[i-start] = createDom(newChn[i], "")
				appendChild(newNode.dom, createDom(newChn[i], ""))
			}
			// newNode.dom.Call("append", doms...)
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
		// chdDom := chNodes.Index(start)
		// var nextNewKey any
		if len(newChn) > start+1 {
			// nextNewKey = newChn[start+1].key
		}
		if mappedOld != nil {
			// oldDom := getDom(mappedOld)
			// if !chdDom.Equal(oldDom) {
			// 	// moveBefore(newNode.dom, nextNewKey, oldChd.key, chdDom, oldDom)
			// }
			reconcile(mappedOld, newChd)
			delete(oldMap, newKey)
		} else {
			// moveBefore(newNode.dom, nextNewKey, oldChd.key, chdDom, createDom(newChd, ""))
		}
		start++
	}

	for _, node := range oldMap {
		// getDom(node).Call("remove")
		node.teardown()
	}
}
