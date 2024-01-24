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

//export createElement
func createElement(tag string, clicks bool) int

//export createTd
func createTd(clicks bool) int

//export createTr
func createTr(clicks bool) int

//export createSpan
func createSpan(clicks bool) int

//export createDiv
func createDiv(clicks bool) int

//export createTable
func createTable(clicks bool) int

//export createTbody
func createTbody(clicks bool) int

//export createH1
func createH1(clicks bool) int

//export createA
func createA(clicks bool) int

//export createButton
func createButton(clicks bool) int

//export createElementNS
func createElementNS(tag string, ns string, clicks bool) int

//export createTextNode
func createTextNode(text string) int

//export appendChild
func appendChild(parent int, child int)

//export replaceWith
func replaceWith(old, new int)

//export moveBefore
func moveBefore(parent int, nextKeyMatch bool, start int, movingDomNode int)

//export mount
func mount(child int, selector string)

//export setStr
func setStr(child int, prop string, val string)

//export setClass
func setClass(child int, val string)

//export setTextContent
func setTextContent(child int, val string)

//export setAriaHidden
func setAriaHidden(child int, val bool)

//export setBool
func setBool(child int, prop string, val bool)

//export removeAttribute
func removeAttribute(child int, attr string)

//export removeNode
func removeNode(node int)

//export disposeNode
func disposeNode(node int)

//export cloneNode
func cloneNode(node int) int

var clickListeners = map[int]func(*MouseEvent){}
var _listener func(*MouseEvent)
var _callClickListener = js.FuncOf(func(js.Value, []js.Value) any {
	_listener(newMouseEvent(global.Get("_GOUI_EVENT")))
	return nil
})

//export callClickListener
func callClickListener(node int) {
	if listener, ok := clickListeners[node]; ok {
		_listener = listener
		_callClickListener.Invoke()
	}
}

var elements = global.Get("_GOUI_ELEMENTS")

func getJsValue(ref int) js.Value {
	return elements.Index(int(ref))
}
