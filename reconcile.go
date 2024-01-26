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
		replaceWith(oldDom, createDom(newNode, oldNode.namespace))
		oldNode.teardown()
	} else if oldNode.render != nil {
		reconcileComponents(oldNode, newNode)
	} else {
		newNode.dom = oldNode.dom
		if oldNode.tag != "" {
			reconcileVdomElems(oldNode, newNode)
		} else if oldNode.textContent != newNode.textContent {
			setData(newNode.dom, newNode.textContent)
		}
	}
}

func reconcileVdomElems(oldNode *Node, newNode *Node) {
	if areDepsEqual(oldNode.memo, newNode.memo) {
		return
	}
	if oldNode.textContent != newNode.textContent {
		setTextContent(newNode.dom, newNode.textContent)
	}
	reconcileAttributes(oldNode, newNode)
	reconcileChildren(oldNode, newNode)
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

func reconcileComponents(oldNode *Node, newNode *Node) {
	newNode.hooks = oldNode.hooks
	if areDepsEqual(oldNode.memo, newNode.memo) {
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
			setBool(ref, name, true)
		}
	}
}

func reconcileAttributes(oldNode *Node, newNode *Node) {
	oldAttrs := oldNode.attrs
	newAttrs := newNode.attrs

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
			setAriaHidden(oldNode.dom, true)
		}
	}
	reconcileBoolAttribute(oldAttrs.Disabled, newAttrs.Disabled, "disabled", oldNode.dom)
	reconcileStringAttribute(oldAttrs.Style, newAttrs.Style, "style", oldNode.dom)
	reconcileStringAttribute(oldAttrs.Class, newAttrs.Class, "class", oldNode.dom)
	reconcileStringAttribute(oldAttrs.ID, newAttrs.ID, "id", oldNode.dom)
	reconcileStringAttribute(oldAttrs.Value, newAttrs.Value, "value", oldNode.dom)
	if oldAttrs.OnClick != newAttrs.OnClick {
		clickListeners[newNode.dom] = newAttrs.OnClick.invoke
	}
}

func callComponentFuncAndReconcile(oldNode *Node, newNode *Node) {
	newElemVdom := callComponentFunc(newNode)
	reconcile(oldNode.virtNode, newElemVdom)
	newNode.virtNode = newElemVdom
}

func reconcileChildren(oldNode *Node, newNode *Node) {
	newChn := newNode.children
	oldChn := oldNode.children
	newLength := len(newChn)
	oldLength := len(oldChn)
	if newLength == 0 && oldLength > 0 {
		setTextContent(newNode.dom, "")
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
		if n.key == "" || n.key == o.key {
			reconcile(o, n)
		} else {
			break
		}
		start++
	}
	if start >= newLength {
		for start < oldLength {
			removeNode(getDom(oldChn[start]))
			oldChn[start].teardown()
			start++
		}
		return
	}

	// suffix
	oldLength--
	newLength--
	for newLength > start && oldLength >= start {
		o := oldChn[oldLength]
		n := newChn[newLength]
		if n.key == "" || n.key == o.key {
			reconcile(o, n)
		} else {
			break
		}
		oldLength--
		newLength--
	}

	oldMap := make(map[string]*Node)
	for i := start; i <= oldLength; i++ {
		oldChd := oldChn[i]
		oldKey := oldChd.key
		noMoreNewChn := false
		if i >= len(newChn) {
			noMoreNewChn = true
		}
		if oldKey != "" && (noMoreNewChn || oldKey != newChn[i].key) {
			oldMap[oldKey] = oldChd
		}
	}

	for start <= newLength {
		if len(oldChn) <= start {
			for i := start; i <= newLength; i++ {
				appendChild(newNode.dom, createDom(newChn[i], newNode.namespace))
			}
			break
		}

		newChd := newChn[start]
		oldChd := oldChn[start]
		if oldChd.key == newChd.key {
			reconcile(oldChd, newChd)
			start++
			continue
		}
		var nextIsCorrect bool
		if len(newChn) > start+1 {
			nextIsCorrect = newChn[start+1].key == oldChd.key
		}
		var oldDom int
		if mappedOld, ok := oldMap[newChd.key]; ok {
			oldDom = getDom(mappedOld)
			reconcile(mappedOld, newChd)
			delete(oldMap, newChd.key)
		} else {
			oldDom = createDom(newChd, newNode.namespace)
		}
		moveBefore(newNode.dom, nextIsCorrect, start, oldDom)
		start++
	}

	for _, node := range oldMap {
		removeNode(getDom(node))
		node.teardown()
	}
}
