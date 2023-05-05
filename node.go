package goui

import (
	"strings"

	"github.com/twharmon/godom"
)

type Node struct {
	dom   *godom.Elem
	tag   string // empty string for text node
	text  string // for text node only
	attrs Attributes

	// TODO: maybe use a map if there are too many
	onClick     *godom.Listener
	onMouseMove *godom.Listener
	onInput     *godom.Listener

	// These fields are only used for component nodes
	id       ID
	vdom     *Node
	props    any
	fn       func(any) *Node
	pc       uintptr
	_states  *store[uintptr, any]
	_effects *store[uintptr, *effectRecord]
	_memos   *store[uintptr, *memoRecord]
}

func (n *Node) AsChildren() []*Node {
	return []*Node{n}
}

func (n *Node) getEffects() *store[uintptr, *effectRecord] {
	if n._effects == nil {
		n._effects = newStore[uintptr, *effectRecord]()
	}
	return n._effects
}

func (n *Node) getMemos() *store[uintptr, *memoRecord] {
	if n._memos == nil {
		n._memos = newStore[uintptr, *memoRecord]()
	}
	return n._memos
}

func (n *Node) getStates() *store[uintptr, any] {
	if n._states == nil {
		n._states = newStore[uintptr, any]()
	}
	return n._states
}

func (n *Node) isComponent() bool {
	return n.fn != nil
}

func (n *Node) isText() bool {
	return n.tag == ""
}

func (n *Node) createDom() *godom.Elem {
	if n.fn != nil {
		rootNode := n.fn(n.props)
		rootDom := rootNode.createDom()
		n.vdom = rootNode
		n.dom = rootDom
		return rootDom
	}
	if n.isText() {
		n.dom = godom.CreateTextElem(n.text)
		return n.dom
	}
	n.dom = godom.Create(n.tag)
	if n.attrs.Disabled {
		n.dom.Attr("disabled", true)
	}
	if len(n.attrs.Class) > 0 {
		n.dom.Classes(strings.Split(n.attrs.Class, " ")...)
	}
	if n.attrs.Style != "" {
		n.dom.Attr("style", n.attrs.Style)
	}
	if n.attrs.Value != "" {
		n.dom.Attr("value", n.attrs.Value)
	}
	if n.attrs.OnClick != nil {
		n.onClick = n.dom.AddMouseEventListener("click", n.attrs.OnClick)
	}
	if n.attrs.OnInput != nil {
		n.onInput = n.dom.AddInputEventListener("input", n.attrs.OnInput)
	}
	if n.attrs.OnMouseMove != nil {
		n.onMouseMove = n.dom.AddMouseEventListener("mousemove", n.attrs.OnMouseMove)
	}
	for _, child := range n.attrs.Children {
		n.dom.AppendChild(child.createDom())
	}
	return n.dom
}

func (n *Node) teardown() {
	if n.fn != nil {
		if n._effects != nil {
			records := n._effects.all()
			for _, record := range records {
				record.teardown()
			}
			for _, child := range n.attrs.Children {
				child.teardown()
			}
			n._effects.clear()
		}
		if n._memos != nil {
			n._memos.clear()
		}
		if n._states != nil {
			n._states.clear()
		}
		componentIDGenerator.release(n.id)
	}
}
