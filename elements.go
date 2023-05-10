package goui

import (
	"fmt"
	"reflect"
	"runtime"
	"runtime/debug"
	"strings"

	"github.com/twharmon/godom"
)

var NoProps any = struct{}{}

var buildInfoPath string

var componentIDGenerator = newIDGenerator()

type Attributes struct {
	Class    string
	Disabled bool
	Style    string
	Value    string
	Children Children

	OnClick     func(*godom.MouseEvent)
	OnMouseMove func(*godom.MouseEvent)
	OnInput     func(*godom.InputEvent)
}

func init() {
	buildInfo, ok := debug.ReadBuildInfo()
	if !ok {
		panic("unable to read build info")
	}
	buildInfoPath = buildInfo.Path + "/"
}

func Text(text string, args ...any) *Node {
	return &Node{
		tag:  "",
		text: fmt.Sprintf(text, args...),
	}
}

func Element(tag string, attrs Attributes) *Node {
	return &Node{
		tag:      tag,
		attrs:    attrs,
		children: renderableToNodes(attrs.Children),
	}
}

func Component[Props any](render func(Props) *Node, props Props) *Node {
	fn := runtime.FuncForPC(uintptr(reflect.ValueOf(render).UnsafePointer()))
	n := &Node{
		props:    props,
		pc:       fn.Entry(),
		name:     strings.TrimPrefix(fn.Name(), buildInfoPath),
		id:       componentIDGenerator.generate(),
		updateCh: make(chan struct{}),
		doneCh:   make(chan struct{}),
	}
	n.render = func() {
		renderWithCurrentNodeLocked(n, render)
		go n.runEffects()
	}
	go func() {
		for {
			select {
			case <-n.updateCh:
				n.update()
			case <-n.doneCh:
				return
			}
		}
	}()
	return n
}
