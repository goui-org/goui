package goui

import "syscall/js"

var global = js.Global()
var document = global.Get("document")

type console struct {
	obj js.Value
}

func (c *console) Log(args ...any) {
	c.obj.Call("log", args...)
}

var Console = &console{obj: global.Get("console")}

var createElementVal = document.Get("createElement")
var createTexNodeFunc = document.Get("createTextNode").Call("bind", document)
var createElementMap = map[string]js.Value{}

func createElement(tag string) js.Value {
	if t, ok := createElementMap[tag]; ok {
		return t.Invoke()
	}
	invoker := createElementVal.Call("bind", document, tag)
	createElementMap[tag] = invoker
	return invoker.Invoke()
}

func createElementNS(tag string, ns string) js.Value {
	return document.Call("createElementNS", ns, tag)
}

func createTextNode(text string) js.Value {
	return createTexNodeFunc.Invoke(text)
}
