package goui

import (
	"fmt"
	"reflect"
	"runtime"

	"github.com/twharmon/godom"
)

type DivAttributes struct {
	Classes  []string
	Style    string
	Text     string
	Children []*Node
	OnClick  func(*godom.MouseEvent)
}

func Div(attrs DivAttributes) *Node {
	return &Node{
		ty:       "div",
		classes:  attrs.Classes,
		style:    attrs.Style,
		text:     attrs.Text,
		children: attrs.Children,
		onClick:  attrs.OnClick,
	}
}

type HeadingAttributes struct {
	Classes  []string
	Style    string
	Text     string
	Children []*Node
	OnClick  func(*godom.MouseEvent)
}

func Heading(level int, attrs HeadingAttributes) *Node {
	return &Node{
		ty:       fmt.Sprintf("h%d", level),
		classes:  attrs.Classes,
		style:    attrs.Style,
		text:     attrs.Text,
		children: attrs.Children,
		onClick:  attrs.OnClick,
	}
}

func Text(text string) *Node {
	return &Node{
		ty:   "",
		text: text,
	}
}

type ButtonAttributes struct {
	Classes  []string
	Disabled bool
	Style    string
	Text     string
	Children []*Node
	OnClick  func(*godom.MouseEvent)
}

func Button(attrs ButtonAttributes) *Node {
	return &Node{
		ty:       "button",
		classes:  attrs.Classes,
		style:    attrs.Style,
		text:     attrs.Text,
		children: attrs.Children,
		onClick:  attrs.OnClick,
		disabled: attrs.Disabled,
	}
}

var NoProps any = struct{}{}

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
