package goui

import (
	"github.com/twharmon/godom"
)

type Node struct {
	dom      *godom.Elem
	ty       string
	classes  []string
	disabled bool
	style    string
	text     string
	children []*Node
	onClick  func(*godom.MouseEvent)

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
	return n.ty == ""
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
	n.dom = godom.Create(n.ty)
	if n.ty == "button" && n.disabled {
		n.dom.Attr("disabled", true)
	}
	// Class    string
	if len(n.classes) > 0 {
		n.dom.Classes(n.classes...)
	}
	// Style    string
	if n.style != "" {
		n.dom.Attr("style", n.style)
	}
	// Text     string
	if n.text != "" {
		n.dom.Text(n.text)
	}
	if n.onClick != nil {
		n.dom.OnClick(n.onClick)
	}
	// Children []Renderer
	for _, child := range n.children {
		el := child.createDom()
		n.dom.AppendChild(el)
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
		for _, child := range n.children {
			child.teardown()
		}
	}
}
