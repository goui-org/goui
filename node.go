package goui

import (
	"strings"
	"sync"

	"github.com/twharmon/godom"
	"github.com/twharmon/goui/utils/concurrentmap"
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
	vdommu   sync.Mutex // TODO: May not need this
	props    any
	fn       func(any) *Node
	pc       uintptr
	_states  *concurrentmap.Map[uintptr, any]
	_effects *concurrentmap.Map[uintptr, *effectRecord]
	_memos   *concurrentmap.Map[uintptr, *memoRecord]
}

func (n *Node) AsChildren() []*Node {
	return []*Node{n}
}

func (n *Node) getEffects() *concurrentmap.Map[uintptr, *effectRecord] {
	if n._effects == nil {
		n._effects = concurrentmap.New[uintptr, *effectRecord]()
	}
	return n._effects
}

func (n *Node) getMemos() *concurrentmap.Map[uintptr, *memoRecord] {
	if n._memos == nil {
		n._memos = concurrentmap.New[uintptr, *memoRecord]()
	}
	return n._memos
}

func (n *Node) getStates() *concurrentmap.Map[uintptr, any] {
	if n._states == nil {
		n._states = concurrentmap.New[uintptr, any]()
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
		n.vdommu.Lock()
		n.vdom = n.fn(n.props) //TODO lock
		rootDom := n.vdom.createDom()
		n.dom = rootDom
		n.vdommu.Unlock()
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
			records := n._effects.AllValues()
			for _, record := range records {
				record.teardown()
			}
			for _, child := range n.attrs.Children {
				child.teardown()
			}
			n._effects.Clear()
		}
		if n._memos != nil {
			n._memos.Clear()
		}
		if n._states != nil {
			n._states.Clear()
		}
		componentIDGenerator.release(n.id)
	}
}
