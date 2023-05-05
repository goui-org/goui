package goui

import "sync"

var currentNodeRef = struct {
	mu   sync.RWMutex
	node *Node
}{}

func UseCurrentComponentID() ID {
	return useCurrentComponent().id
}

func useCurrentComponent() *Node {
	currentNodeRef.mu.RLock()
	node := currentNodeRef.node
	currentNodeRef.mu.RUnlock()
	return node
}

func assignCurrentComponent(n *Node) {
	currentNodeRef.mu.Lock()
	currentNodeRef.node = n
	currentNodeRef.mu.Unlock()
}
