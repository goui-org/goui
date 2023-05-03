package goui

import (
	"github.com/twharmon/godom"
)

type Node struct {
	dom   *godom.Elem
	tag   string // empty string for text node
	text  string // for text node only
	attrs Attributes

	// These fields are only used for component nodes
	vdom    *Node
	props   any
	fn      func(any) *Node
	id      uintptr
	state   *store[uintptr, any]
	effects *store[uintptr, *effectRecord]
}

func (n *Node) Render() *Node {
	return n
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
		n.dom.Attr("class", n.attrs.Class)
	}
	if n.attrs.Style != "" {
		n.dom.Attr("style", n.attrs.Style)
	}
	if n.attrs.Value != "" {
		n.dom.Attr("value", n.attrs.Value)
	}
	if n.attrs.OnClick != nil {
		n.dom.OnClick(n.attrs.OnClick)
	}
	if n.attrs.OnInput != nil {
		n.dom.OnInput(n.attrs.OnInput)
	}
	for _, child := range n.attrs.Children {
		n.dom.AppendChild(child.createDom())
	}
	return n.dom
}

func (n *Node) teardown() {
	if n.fn != nil {
		records := n.effects.all()
		for _, record := range records {
			record.teardown()
		}
		components.delete(n.id)
		n.state = nil
		for _, child := range n.attrs.Children {
			child.teardown()
		}
	}
}
