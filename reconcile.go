package goui

import (
	"fmt"
	"reflect"
	"strings"
)

func reconcile(old *Node, new *Node) {
	if old.ty != new.ty || old.id != new.id {
		newDom := new.createDom()
		old.teardown()
		old.dom.Replace(newDom)
		return
	}
	if old.isComponent() {
		// both old and new are same type components
		reconcileVdomComponents(old, new)
		return
	}
	// both old and new are plain dom nodes
	reconcileVdomNodes(old, new)
}

func reconcileVdomComponents(old *Node, new *Node) {
	if fmt.Sprintf("%v", old.props) == fmt.Sprintf("%v", new.props) {
		// new.instance = old.instance
		new.vdom = old.vdom
		new.dom = old.dom
		return
	}
	// old.props = new.props
	new.vdom = new.fn(new.props)
	new.dom = old.dom
	new.vdom.dom = old.dom
	reconcile(old.vdom, new.vdom)
	// old.instance.setProps(new.props)
	// new.instance = old.instance
	// new.vdom = render(new.instance)
	// new.dom = old.dom
	// new.vdom.dom = old.dom
	// reconcile(old.vdom, new.vdom)
}

func reconcileVdomNodes(old *Node, new *Node) {
	if old.ty != new.ty {
		newDom := new.createDom()
		old.dom.Replace(newDom)
		return
	}
	if old.text != new.text {
		old.dom.Text(new.text)
	}
	if strings.Join(old.classes, ",") != strings.Join(new.classes, ",") {
		old.dom.Classes(new.classes...)
	}
	if old.style != new.style {
		old.dom.Attr("style", new.style)
	}
	if old.disabled != new.disabled {
		old.dom.Attr("disabled", new.disabled)
	}
	if new.onClick != nil {
		if old.onClick != nil {
			oldFnPtr := reflect.ValueOf(old.onClick).UnsafePointer()
			newFnPtr := reflect.ValueOf(new.onClick).UnsafePointer()
			if oldFnPtr != newFnPtr {
				old.dom.OnClick(new.onClick)
			}
		} else {
			old.dom.OnClick(new.onClick)
		}
	}
	if len(old.children) > len(new.children) {
		for i := len(new.children); i < len(old.children); i++ {
			old.children[i].teardown()
			old.dom.RemoveChild(old.children[i].dom)
		}
	} else if len(old.children) < len(new.children) {
		for i := len(old.children); i < len(new.children); i++ {
			old.dom.AppendChild(new.children[i].createDom())
		}
	}
	new.dom = old.dom
	for i := 0; i < min(len(old.children), len(new.children)); i++ {
		newNode := new.children[i]
		oldNode := old.children[i]
		reconcile(oldNode, newNode)
	}
}
