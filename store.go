package goui

// import (
// 	"syscall/js"
// )

// type elemStore struct {
// 	vals map[string]*stack
// }

// var store = &elemStore{
// 	vals: make(map[string]*stack),
// }

// func (es *elemStore) put(e *Elem) {
// 	if e.tag == "" {
// 		e.dom.Set("data", nil)
// 	} else {
// 		for name, fn := range e.listeners {
// 			fn.Release()
// 			e.dom.Set(name, nil)
// 		}
// 		attrs := e.attrs.(*Attributes)
// 		if attrs.Disabled {
// 			e.dom.Set("disabled", false)
// 		}
// 		if attrs.Class != "" {
// 			e.dom.Call("removeAttribute", "class")
// 		}
// 		if attrs.Style != "" {
// 			e.dom.Call("removeAttribute", "style")
// 		}
// 		if attrs.ID != "" {
// 			e.dom.Call("removeAttribute", "id")
// 		}
// 		if attrs.AriaHidden {
// 			e.dom.Call("removeAttribute", "aria-hidden")
// 		}
// 		if attrs.Value != "" {
// 			e.dom.Set("value", nil)
// 		}
// 	}
// 	stack := es.vals[e.tag]
// 	if stack == nil {
// 		stack = newJsValueStack()
// 		es.vals[e.tag] = stack

// 	}
// 	stack.Push(e.dom)
// }

// func (es *elemStore) getElement(tag string) js.Value {
// 	vals := es.vals[tag]
// 	if vals != nil && vals.Len() > 0 {
// 		return vals.Pop()
// 	}
// 	e := createElement(tag)
// 	return e
// }

// func (es *elemStore) getTextNode(text string) js.Value {
// 	vals := es.vals[""]
// 	if vals != nil && vals.Len() > 0 {
// 		val := vals.Pop()
// 		val.Set("data", text)
// 		return val
// 	}
// 	return createTextNode(text)
// }
