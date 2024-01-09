package goui

import "syscall/js"

// MouseEvent .
type MouseEvent struct {
	val js.Value
}

func newMouseEvent(val js.Value) *MouseEvent {
	return &MouseEvent{
		val: val,
	}
}

// PreventDefault .
func (e *MouseEvent) PreventDefault() {
	e.val.Call("preventDefault")
}

// StopPropogation .
func (e *MouseEvent) StopPropogation() {
	e.val.Call("stopPropogation")
}

// OffsetX .
func (e *MouseEvent) OffsetX() int {
	return e.val.Get("offsetX").Int()
}

// OffsetY .
func (e *MouseEvent) OffsetY() int {
	return e.val.Get("offsetY").Int()
}

// ClientX .
func (e *MouseEvent) ClientX() int {
	return e.val.Get("clientX").Int()
}

// ClientY .
func (e *MouseEvent) ClientY() int {
	return e.val.Get("clientY").Int()
}

// ShiftKey .
func (e *MouseEvent) ShiftKey() bool {
	return e.val.Get("shiftKey").Bool()
}

// AltKey .
func (e *MouseEvent) AltKey() bool {
	return e.val.Get("altKey").Bool()
}

// CtrlKey .
func (e *MouseEvent) CtrlKey() bool {
	return e.val.Get("ctrlKey").Bool()
}
