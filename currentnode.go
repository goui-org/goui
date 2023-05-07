package goui

import (
	"sync"
	"time"

	"github.com/twharmon/godom"
)

var currentNodeRef = struct {
	mu   sync.Mutex
	node *Node
}{}

func UseCurrentComponentID() ID {
	return useCurrentComponent().id
}

func useCurrentComponent() *Node {
	return currentNodeRef.node
}

func renderWithCurrentNodeLocked[T any](n *Node, render func(T) *Node) *Node {
	currentNodeRef.mu.Lock()
	currentNodeRef.node = n
	t := time.Now()
	n.vdom = render(n.props.(T)) // hooks are protected to use useCurrentComponent()
	if dur := time.Since(t); dur > time.Millisecond*5 {
		godom.Console.Warn("[GOUI] Warning: %s took %s to render", n.name, dur)
	}
	currentNodeRef.mu.Unlock()
	return n.vdom
}
