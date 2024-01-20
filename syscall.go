package goui

import "syscall/js"

var global = js.Global()
var document = global.Get("document")
var console = global.Get("console")

func Log(args ...any) {
	console.Call("log", args...)
}
