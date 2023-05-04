package goui

import (
	"fmt"
	"reflect"
	"runtime"

	"github.com/twharmon/godom"
)

func Text(text string, args ...any) *Node {
	return &Node{
		tag:  "",
		text: fmt.Sprintf(text, args...),
	}
}

var NoProps any = struct{}{}

type Attributes struct {
	Class    string
	Disabled bool
	Style    string
	Value    string
	Children []*Node

	OnClick     func(*godom.MouseEvent)
	OnMouseMove func(*godom.MouseEvent)
	OnInput     func(*godom.InputEvent)
}

func Element(tag string, attrs Attributes) *Node {
	return &Node{
		tag:   tag,
		attrs: attrs,
	}
}

func Component[Props any](fn func(Props) *Node, props Props) *Node {
	pc := runtime.FuncForPC(uintptr(reflect.ValueOf(fn).UnsafePointer())).Entry()
	n := &Node{
		props: props,
		pc:    pc,
	}
	n.fn = func(p any) *Node {
		prev := getCurrentNode()
		assignCurrentNode(n)
		node := fn(p.(Props))
		assignCurrentNode(prev)
		return node
	}
	return n
}
