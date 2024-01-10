package goui

import "syscall/js"

type domNode struct {
	val js.Value

	removeClass      js.Value
	removeAriaHidden js.Value
}

// 	e.dom.Set("disabled", false)
// 	e.dom.Call("removeAttribute", "class")
// 	e.dom.Call("removeAttribute", "style")
// 	e.dom.Call("removeAttribute", "id")
// 	e.dom.Call("removeAttribute", "aria-hidden")
// 	e.dom.Set("value", nil)

func newDomNode(val js.Value) *domNode {
	removeAttr := val.Get("removeAttribute")
	return &domNode{
		val: val,
		// clearText: val.Get(),
		removeClass:      removeAttr.Call("bind", val, "class"),
		removeAriaHidden: removeAttr.Call("bind", val, "aria-hidden"),
	}
}

func (d *domNode) resetElement() {
	d.removeClass.Invoke()
	d.removeAriaHidden.Invoke()
}

func (d *domNode) resetTextNode() {
	d.val.Set("data", nil)
}
