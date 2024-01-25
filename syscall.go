package goui

import (
	"syscall/js"
)

var global = js.Global()
var document = global.Get("document")
var console = global.Get("console")

func Log(args ...any) {
	console.Call("log", args...)
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

//export setData
func setData(child int, val string)

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

// func init() {
// 	var memStats runtime.MemStats
// 	go func() {
// 		for range time.NewTicker(time.Second).C {
// 			runtime.ReadMemStats(&memStats)
// 			fmt.Printf("HeapInuse=%d;Frees=%d\n", memStats.HeapInuse/1e6, memStats.Frees)
// 		}
// 	}()
// }

var _listener func(*MouseEvent)
var _callClickListener = js.FuncOf(func(this js.Value, args []js.Value) any {
	_listener(newMouseEvent(args[0]))
	return nil
})

//export callClickListener
func callClickListener(node int) {
	if listener, ok := clickListeners[node]; ok {
		// println("listener called")
		// listener(newMouseEvent(global.Get("_GOUI_EVENT")))
		// runtime.GC()
		_listener = listener
		_callClickListener.Invoke(global.Get("_GOUI_EVENT"))
	}
}

var elements = global.Get("_GOUI_ELEMENTS")

func getJsValue(ref int) js.Value {
	return elements.Index(int(ref))
}
