package goui

import "sync"

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

func renderWithCurrentNodeLocked[T any](n *Node, fn func(T) *Node, props T) *Node {
	currentNodeRef.mu.Lock()
	currentNodeRef.node = n
	// n.vdommu.Lock()
	n.vdom = fn(props) // hooks are protected
	// n.vdommu.Unlock()
	currentNodeRef.mu.Unlock()
	return n.vdom
}
