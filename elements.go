package goui

import (
	"reflect"
	"runtime"

	"github.com/twharmon/godom"
)

func Text(text string) *Node {
	return &Node{
		tag:  "",
		text: text,
	}
}

var NoProps any = struct{}{}

type Attributes struct {
	Class    string
	Disabled bool
	Style    string
	Value    string
	Children []*Node

	OnClick func(*godom.MouseEvent)
	OnInput func(*godom.InputEvent)
}

func Element(tag string, attrs Attributes) *Node {
	return &Node{
		tag:   tag,
		attrs: attrs,
	}
}

func Component[Props any](fn func(Props) *Node, props Props) *Node {
	id := runtime.FuncForPC(uintptr(reflect.ValueOf(fn).UnsafePointer())).Entry()
	n := &Node{
		props:   props,
		fn:      func(p any) *Node { return fn(p.(Props)) },
		state:   newStore[uintptr, any](),
		effects: newStore[uintptr, *effectRecord](),
		id:      id,
	}
	components.set(id, n)
	return n
}
