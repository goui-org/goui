package goui

import (
	"strings"

	"github.com/twharmon/godom"
	"github.com/twharmon/goui/utils/equalityutil"
	"github.com/twharmon/goui/utils/mathutil"
)

func reconcile(old *Node, new *Node) {
	if old.tag != new.tag || old.pc != new.pc {
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
	reconcileVdomNodes(old, new) // sendint old.dom == nil here
}

func reconcileVdomComponents(old *Node, new *Node) {
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

	if old.attrs.Class != new.attrs.Class {
		old.dom.Classes(strings.Split(new.attrs.Class, " ")...)
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
		elem.Attr(name, newAttr)
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
	commonCnt := mathutil.Min(len(old.attrs.Children), len(new.attrs.Children))
	for i := 0; i < commonCnt; i++ {
		reconcile(old.attrs.Children[i], new.attrs.Children[i])
	}
}
