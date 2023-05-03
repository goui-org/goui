package goui

import (
	"fmt"
	"reflect"

	"github.com/twharmon/godom"
)

func reconcile(old *Node, new *Node) {
	if old.tag != new.tag || old.id != new.id {
		newDom := new.createDom()
		old.teardown()
		old.dom.ReplaceWith(newDom)
		return
	}
	if old.isComponent() {
		// both old and new are same type components
		reconcileVdomComponents(old, new)
		return
	}
	// both old and new are plain dom nodes of same type
	reconcileVdomNodes(old, new)
}

func reconcileVdomComponents(old *Node, new *Node) {
	if fmt.Sprintf("%v", old.props) == fmt.Sprintf("%v", new.props) {
		new.vdom = old.vdom
		new.dom = old.dom
		return
	}
	new.vdom = new.fn(new.props)
	new.dom = old.dom
	reconcile(old.vdom, new.vdom)
}

func reconcileVdomNodes(old *Node, new *Node) {
	if old.isText() {
		if old.text != new.text {
			old.dom.Text(new.text)
		}
		new.dom = old.dom
		return
	}

	// attributes
	reconcileAttribute(old.attrs.Class, new.attrs.Class, "class", old.dom)
	reconcileAttribute(old.attrs.Style, new.attrs.Style, "style", old.dom)
	reconcileAttribute(old.attrs.Disabled, new.attrs.Disabled, "disabled", old.dom)
	reconcileAttribute(old.attrs.Value, new.attrs.Value, "value", old.dom)

	// listeners
	reconcileListener(old.attrs.OnClick, new.attrs.OnClick, old.dom.OnClick)
	reconcileListener(old.attrs.OnInput, new.attrs.OnInput, old.dom.OnInput)

	new.dom = old.dom
	reconcileChildren(old, new)
}

func reconcileAttribute[T comparable](oldAttr T, newAttr T, name string, elem *godom.Elem) {
	if oldAttr != newAttr {
		elem.Attr(name, newAttr)
	}
}

func reconcileListener[T any](oldFn func(T), newFn func(T), setter func(func(T)) *godom.Elem) {
	if newFn != nil {
		if oldFn != nil {
			oldFnPtr := reflect.ValueOf(oldFn).UnsafePointer()
			newFnPtr := reflect.ValueOf(newFn).UnsafePointer()
			if oldFnPtr != newFnPtr {
				setter(newFn)
			}
		} else {
			setter(newFn)
		}
	} else {
		if oldFn != nil {
			var t func(T)
			setter(t)
		}
	}
}

func reconcileChildren(old *Node, new *Node) {
	if len(old.attrs.Children) > len(new.attrs.Children) {
		// previous dom has more children, teardown and remove them
		for i := len(new.attrs.Children); i < len(old.attrs.Children); i++ {
			old.attrs.Children[i].teardown()
			old.dom.RemoveChild(old.attrs.Children[i].dom)
		}
	} else if len(old.attrs.Children) < len(new.attrs.Children) {
		// previous dom has fewer children, create the new ones
		for i := len(old.attrs.Children); i < len(new.attrs.Children); i++ {
			old.dom.AppendChild(new.attrs.Children[i].createDom())
		}
	}
	commonCnt := min(len(old.attrs.Children), len(new.attrs.Children))
	for i := 0; i < commonCnt; i++ {
		reconcile(old.attrs.Children[i], new.attrs.Children[i])
	}
}
