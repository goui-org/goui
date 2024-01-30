package goui

import (
	"reflect"
	"strconv"
	"syscall/js"
)

type NoProps any

type Node struct {
	tag         string
	namespace   string
	ptr         uintptr
	render      func() *Node
	key         string
	attrs       *Attributes
	textContent string
	ref         *Ref[js.Value]
	refs        []int
	dom         int
	unmounted   bool
	children    []*Node

	virtNode    *Node
	queue       []*Node
	hooks       []any
	hooksCursor int
	memo        Deps
}

func (n *Node) teardown() {
	if n.virtNode != nil {
		n.unmounted = true
		n.queue = nil
		for _, h := range n.hooks {
			if effect, ok := h.(*effectRecord); ok {
				if effect.teardown != nil {
					effect.teardown()
				}
			}
		}
		n.virtNode.teardown()
		return
	}
	if n.ref != nil {
		n.ref.Value = js.Undefined()
	}
	for _, ch := range n.children {
		ch.teardown()
	}
	if n.attrs != nil && n.attrs.OnClick != nil {
		delete(clickListeners, n.dom)
	}
	disposeNode(n.dom)
	for _, n := range n.refs {
		disposeNode(n)
	}
	n.refs = n.refs[:0]
	n.dom = 0
}

type Keyer interface {
	Key() string
}

type Memoer interface {
	Memo() Deps
}

func Component[T any](ty func(T) *Node, props T) *Node {
	n := &Node{
		ptr:    uintptr(reflect.ValueOf(ty).UnsafePointer()),
		render: func() *Node { return ty(props) },
	}
	if keyer, ok := any(props).(Keyer); ok {
		n.key = keyer.Key()
	}
	if memoer, ok := any(props).(Memoer); ok {
		n.memo = memoer.Memo()
	}
	return n
}

var currentNode *Node

func callComponentFunc(node *Node) *Node {
	prev := currentNode
	currentNode = node
	node.hooksCursor = 0
	vd := node.render()
	if vd == nil {
		vd = &Node{}
	}
	currentNode = prev
	return vd
}

type Attributes struct {
	ID       string
	Class    string
	Disabled bool
	Style    string
	Value    string
	Key      string

	// Slot must be string, int, *Node, []*Node, func(NoProps) *Node, or nil
	Slot any
	Type string
	Ref  *Ref[js.Value]

	AriaHidden bool

	// Common UIEvents: https://developer.mozilla.org/en-US/docs/Web/API/UI_Events
	// All Events:      https://developer.mozilla.org/en-US/docs/Web/API/Event
	//
	// MouseEvent: click, dblclick, mouseup, mousedown
	// InputEvent: input, beforeinput
	// KeyboardEvent: keydown, keypress, keyup
	// CompositionEvent: compositionstart, compositionend, compositionupdate
	// WheelEvent: wheel
	// FocusEvent: focus, blur, focusin, and focusout

	OnClick *Callback[func(*MouseEvent)]
	// OnMouseMove *Callback[func(*MouseEvent)]
	// OnInput     *Callback[func(*InputEvent)]
}

func Element(tag string, attrs *Attributes) *Node {
	n := &Node{
		tag:   tag,
		attrs: attrs,
		key:   attrs.Key,
		ref:   attrs.Ref,
	}
	if attrs.Slot != nil {
		switch slot := attrs.Slot.(type) {
		case string:
			n.textContent = slot
		case int:
			n.textContent = strconv.Itoa(slot)
		case *Node:
			n.children = []*Node{slot}
		case func(NoProps) *Node:
			n.children = []*Node{Component(slot, nil)}
		case []*Node:
			n.children = slot
		case []any:
			n.children = make([]*Node, len(slot))
			for i := 0; i < len(slot); i++ {
				n.children[i] = makeNode(slot[i])
			}
		case []func(NoProps) *Node:
			n.children = make([]*Node, len(slot))
			for i := 0; i < len(slot); i++ {
				n.children[i] = Component(slot[i], nil)
			}
		case []string:
			n.children = make([]*Node, len(slot))
			for i := 0; i < len(slot); i++ {
				n.children[i] = text(slot[i])
			}
		case []int:
			n.children = make([]*Node, len(slot))
			for i := 0; i < len(slot); i++ {
				n.children[i] = text(strconv.Itoa(slot[i]))
			}
		}
	}
	return n
}

func text(content string) *Node {
	return &Node{
		textContent: content,
	}
}

func makeNode(v any) *Node {
	switch chn := v.(type) {
	case string:
		return text(chn)
	case int:
		return text(strconv.Itoa(chn))
	case *Node:
		return chn
	case func(NoProps) *Node:
		return Component(chn, nil)
	}
	return nil
}

func createDom(node *Node, ns string) int {
	if node.dom != 0 {
		node.refs = append(node.refs, node.dom)
		node.dom = cloneNode(node.dom)
		return node.dom
	}
	if node.tag != "" {
		clicks := node.attrs.OnClick != nil
		if node.tag == "svg" {
			ns = "http://www.w3.org/2000/svg"
		} else if node.tag == "math" {
			ns = "http://www.w3.org/1998/Math/MathML"
		}
		switch node.tag {
		case "tr":
			node.dom = createTr(clicks)
		case "span":
			node.dom = createSpan(clicks)
		case "td":
			node.dom = createTd(clicks)
		case "a":
			node.dom = createA(clicks)
		case "h1":
			node.dom = createH1(clicks)
		case "div":
			node.dom = createDiv(clicks)
		case "table":
			node.dom = createTable(clicks)
		case "tbody":
			node.dom = createTbody(clicks)
		case "button":
			node.dom = createButton(clicks)
		default:
			if ns == "" {
				node.dom = createElement(node.tag, clicks)
			} else {
				node.dom = createElementNS(node.tag, ns, clicks)
			}
		}
		if node.ref != nil {
			node.ref.Value = getJsValue(node.dom)
		}
		if node.attrs.Disabled {
			setBool(node.dom, "disabled", true)
		}
		if node.attrs.Class != "" {
			setClass(node.dom, node.attrs.Class)
		}
		if node.attrs.Style != "" {
			setStr(node.dom, "style", node.attrs.Style)
		}
		if node.attrs.ID != "" {
			setStr(node.dom, "id", node.attrs.ID)
		}
		if node.attrs.AriaHidden {
			setAriaHidden(node.dom, true)
		}
		if node.attrs.Value != "" {
			setStr(node.dom, "value", node.attrs.Value)
		}
		if node.textContent != "" {
			setTextContent(node.dom, node.textContent)
		}
		if clicks {
			clickListeners[node.dom] = node.attrs.OnClick.invoke
		}
		node.namespace = ns
		for _, child := range node.children {
			appendChild(node.dom, createDom(child, ns))
		}
	} else if node.render != nil {
		node.virtNode = callComponentFunc(node)
		return createDom(node.virtNode, ns)
	} else {
		node.dom = createTextNode(node.textContent)
		return node.dom
	}
	return node.dom
}
