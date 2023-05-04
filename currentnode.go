package goui

import "sync"

var currentNodeRef = struct {
	mu   sync.RWMutex
	node *Node
}{}

func getCurrentNode() *Node {
	currentNodeRef.mu.RLock()
	node := currentNodeRef.node
	currentNodeRef.mu.RUnlock()
	return node
}

func assignCurrentNode(n *Node) {
	currentNodeRef.mu.Lock()
	currentNodeRef.node = n
	currentNodeRef.mu.Unlock()
}
