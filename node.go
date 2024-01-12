package goui

import (
	"reflect"
	"syscall/js"
)

type NoProps any

type Children []*Node

type Node struct {
	tag       string
	ptr       uintptr
	render    func() *Node
	key       any
	attrs     any
	ref       *Ref[js.Value]
	dom       js.Value
	unmounted bool
	listeners map[string]js.Func

	virtNode    *Node
	queue       []*Node
	hooks       []any
	hooksCursor int
	memo        []any
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
	if attrs, ok := n.attrs.(*Attributes); ok {
		for _, ch := range attrs.Children {
			ch.teardown()
		}
	}
}

type Keyer interface {
	Key() any
}

type Memoer interface {
	Memo() Deps
}

func Component[T any](ty func(T) *Node, props T) *Node {
	fn := uintptr(reflect.ValueOf(ty).UnsafePointer())
	n := &Node{
		ptr:    fn,
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

func Text(content string) *Node {
	return &Node{
		attrs: content,
	}
}

var currentElem *Node

func callComponentFunc(node *Node) *Node {
	prev := currentElem
	currentElem = node
	node.hooksCursor = 0
	vd := node.render()
	if vd == nil {
		vd = &Node{}
	}
	currentElem = prev
	return vd
}

type Attributes struct {
	ID       string
	Class    string
	Disabled bool
	Style    string
	Value    string
	Children Children
	Key      string
	Type     string

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
	return &Node{
		tag:   tag,
		attrs: attrs,
		key:   attrs.Key,
	}
}

var namespacePrefix = "http://www.w3.org/"
var svgNamespace = namespacePrefix + "2000/svg"
var mathNamespace = namespacePrefix + "1998/Math/MathML"

func createDom(node *Node, ns string) js.Value {
	if node.tag != "" {
		if node.tag == "svg" {
			ns = svgNamespace
			node.dom = createElementNS(node.tag, ns)
		} else if node.tag == "math" {
			ns = mathNamespace
			node.dom = createElementNS(node.tag, ns)
		} else {
			node.dom = createElement(node.tag)
		}
		if node.ref != nil {
			node.ref.Value = node.dom
		}
		attrs := node.attrs.(*Attributes)
		if attrs.Disabled {
			node.dom.Set("disabled", true)
		}
		if attrs.Class != "" {
			node.dom.Set("className", attrs.Class)
		}
		if attrs.Style != "" {
			node.dom.Set("style", attrs.Style)
		}
		if attrs.ID != "" {
			node.dom.Set("id", attrs.ID)
		}
		if attrs.AriaHidden {
			node.dom.Set("ariaHidden", true)
		}
		if attrs.Value != "" {
			node.dom.Set("value", attrs.Value)
		}
		if attrs.OnClick != nil {
			node.setEventListener("onclick", func(_ js.Value, args []js.Value) any {
				attrs.OnClick.invoke(newMouseEvent(args[0]))
				return nil
			})
		}
		doms := make([]any, len(attrs.Children))
		for i, child := range attrs.Children {
			doms[i] = createDom(child, ns)
		}
		node.dom.Call("append", doms...)
	} else if node.render != nil {
		node.virtNode = callComponentFunc(node)
		return createDom(node.virtNode, ns)
	} else {
		node.dom = createTextNode(node.attrs.(string))
		return node.dom
	}
	return node.dom
}

func (n *Node) setEventListener(name string, fn func(js.Value, []js.Value) any) {
	if n.listeners == nil {
		n.listeners = make(map[string]js.Func)
	} else if fn := n.listeners[name]; fn.Truthy() {
		fn.Release()
	}
	wrapper := js.FuncOf(fn)
	n.listeners[name] = wrapper
	n.dom.Set(name, wrapper)
}
