package goui

import (
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
		oldDom := getDom(oldNode)
		// TODO: namespace?
		replaceWith(oldDom, createDom(newNode, ""))
		oldNode.teardown()
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
		setStr(newNode.dom, "data", newNode.attrs.(string))
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

func reconcileStringAttribute(oldAttr string, newAttr string, name string, ref int) {
	if oldAttr != newAttr {
		if newAttr == "" {
			removeAttribute(ref, name)
		} else {
			setStr(ref, name, newAttr)
		}
	}
}

func reconcileBoolAttribute(oldAttr bool, newAttr bool, name string, ref int) {
	if oldAttr != newAttr {
		if !newAttr {
			removeAttribute(ref, name)
		} else {
			setBool(ref, name, 1)
		}
	}
}

func reconcileAttributes(oldNode *Node, newNode *Node) {
	oldAttrs := oldNode.attrs.(*Attributes)
	newAttrs := newNode.attrs.(*Attributes)

	if oldAttrs.Class != newAttrs.Class {
		if newAttrs.Class == "" {
			removeAttribute(oldNode.dom, "class")
		} else {
			setStr(oldNode.dom, "className", newAttrs.Class)
		}
	}

	if oldAttrs.AriaHidden != newAttrs.AriaHidden {
		if !newAttrs.AriaHidden {
			removeAttribute(oldNode.dom, "aria-hidden")
		} else {
			setBool(oldNode.dom, "ariaHidden", 1)
		}
	}
	reconcileBoolAttribute(oldAttrs.Disabled, newAttrs.Disabled, "disabled", oldNode.dom)
	reconcileStringAttribute(oldAttrs.Style, newAttrs.Style, "style", oldNode.dom)
	reconcileStringAttribute(oldAttrs.Class, newAttrs.Class, "class", oldNode.dom)
	reconcileStringAttribute(oldAttrs.ID, newAttrs.ID, "id", oldNode.dom)
	reconcileStringAttribute(oldAttrs.Value, newAttrs.Value, "value", oldNode.dom)
	if oldAttrs.OnClick != newAttrs.OnClick {
		oldNode.setEventListener("onclick", func(_ js.Value, args []js.Value) any {
			newAttrs.OnClick.invoke(newMouseEvent(args[0]))
			return nil
		})
	}
}

func reconcileReference(oldNode *Node, newNode *Node) {
	if newNode.ref != nil {
		if oldNode.ref != nil {
			newNode.ref.Value = oldNode.ref.Value
		} else {
			newNode.ref.Value = getJsValue(newNode.dom)
		}
	} else if oldNode.ref != nil {
		oldNode.ref.Value = js.Undefined()
	}
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
	chNodes := getJsValue(newNode.dom).Get("childNodes")
	for start <= newLength {
		newChd := newChn[start]
		if len(oldChn) <= start {
			for i := start; i <= newLength; i++ {
				appendChild(newNode.dom, createDom(newChn[i], ""))
			}
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
			oldDomJsVal := getJsValue(oldDom)
			if !chdDom.Equal(oldDomJsVal) {
				moveBefore(getJsValue(newNode.dom), nextNewKey, oldChd.key, chdDom, oldDomJsVal)
			}
			reconcile(mappedOld, newChd)
			delete(oldMap, newKey)
		} else {
			moveBefore(getJsValue(newNode.dom), nextNewKey, oldChd.key, chdDom, getJsValue(createDom(newChd, "")))
		}
		start++
	}

	for _, node := range oldMap {
		// removeNode(getDom(node))
		node.teardown()
	}
}
