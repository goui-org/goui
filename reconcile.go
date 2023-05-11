package goui

import (
	"fmt"

	"github.com/twharmon/godom"
	"github.com/twharmon/goui/utils/equalityutil"
	"github.com/twharmon/goui/utils/mathutil"
)

func reconcile(old *Node, new *Node) {
	if old.tag != new.tag || old.pc != new.pc {
		new.createDom()
		old.teardown()
		old.dom.ReplaceWith(new.dom)
		return
	}
	if old.isComponent() {
		// both old and new are same type components
		reconcileVdomComponents(old, new)
		return
	}
	// both old and new are plain dom nodes of same type
	reconcileVdomNodes(old, new) // sendint old.dom == nil here
}

func reconcileVdomComponents(old *Node, new *Node) {
	old.done()
	new._effects = old._effects
	new._memos = old._memos
	new._states = old._states
	new.onClick = old.onClick
	new.onInput = old.onInput
	new.onMouseMove = old.onMouseMove
	if equalityutil.DeepEqual(old.props, new.props) {
		new.vdom = old.vdom
		new.dom = old.dom
		return
	}
	new.dom = old.dom
	new.render()
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

	if old.attrs.Class != new.attrs.Class {
		old.dom.Class(new.attrs.Class)
	}

	// attributes
	// reconcileAttribute(old.attrs.Class, new.attrs.Class, "class", old.dom)
	reconcileAttribute(old.attrs.Style, new.attrs.Style, "style", old.dom)
	reconcileAttribute(old.attrs.Disabled, new.attrs.Disabled, "disabled", old.dom)
	reconcileAttribute(old.attrs.Value, new.attrs.Value, "value", old.dom)

	// listeners
	if !equalityutil.DeepEqual(old.attrs.OnClick, new.attrs.OnClick) {
		if old.onClick != nil {
			old.onClick.Remove()
		}
		if new.attrs.OnClick != nil {
			old.dom.AddMouseEventListener("click", new.attrs.OnClick)
		}
	}
	if !equalityutil.DeepEqual(old.attrs.OnInput, new.attrs.OnInput) {
		if old.onInput != nil {
			old.onInput.Remove()
		}
		if new.attrs.OnInput != nil {
			old.dom.AddInputEventListener("input", new.attrs.OnInput)
		}
	}
	if !equalityutil.DeepEqual(old.attrs.OnMouseMove, new.attrs.OnMouseMove) {
		if old.onMouseMove != nil {
			old.onMouseMove.Remove()
		}
		if new.attrs.OnMouseMove != nil {
			old.dom.AddMouseEventListener("mousemove", new.attrs.OnMouseMove)
		}
	}

	new.dom = old.dom
	reconcileChildren(old, new)
}

func reconcileAttribute[T comparable](oldAttr T, newAttr T, name string, elem *godom.Elem) {
	if oldAttr != newAttr {
		var t T
		if newAttr == t {
			fmt.Println("remove attr")
			elem.RemoveAttr(name)
		} else {
			fmt.Println("set attr", newAttr)
			elem.Attr(name, newAttr)
		}
	}
}

func reconcileChildren(old *Node, new *Node) {
	if len(old.children) > len(new.children) {
		// previous dom has more children, teardown and remove them
		for i := len(new.children); i < len(old.children); i++ {
			old.children[i].teardown()
			old.dom.RemoveChild(old.children[i].dom)
		}
	} else if len(old.children) < len(new.children) {
		// previous dom has fewer children, create the new ones
		for i := len(old.children); i < len(new.children); i++ {
			new.children[i].createDom()
			old.dom.AppendChild(new.children[i].dom)
		}
	}
	commonCnt := mathutil.Min(len(old.children), len(new.children))
	for i := 0; i < commonCnt; i++ {
		reconcile(old.children[i], new.children[i])
	}
}
