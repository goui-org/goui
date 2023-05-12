package goui

import (
	"sync"
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
	n.vdom = render(n.props.(T))
	currentNodeRef.mu.Unlock()
	return n.vdom
}
