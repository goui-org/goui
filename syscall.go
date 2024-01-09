package goui

import "syscall/js"

var global = js.Global()
var document = global.Get("document")
var console = global.Get("console")

var Console = struct {
	Log func(args ...any)
}{
	Log: func(args ...any) {
		console.Call("log", args...)
	},
}

func createElement(tag string) js.Value {
	return document.Call("createElement", tag)
}

func createElementNS(tag string, ns string) js.Value {
	return document.Call("createElementNS", tag, ns)
}

func createTextNode(text string) js.Value {
	return document.Call("createTextNode", text)
}
