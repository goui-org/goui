package goui

import (
	"crypto/rand"
	"sync"
)

type ID [8]byte

type idGenerator struct {
	mu    sync.Mutex
	taken map[ID]bool
}

func newIDGenerator() *idGenerator {
	return &idGenerator{
		taken: make(map[ID]bool),
	}
}

func (g *idGenerator) generate() ID {
	g.mu.Lock()
	var id ID
	rand.Read(id[:])
	for g.taken[id] {
		rand.Read(id[:])
	}
	g.taken[id] = true
	g.mu.Unlock()
	return id
}

func (g *idGenerator) release(id ID) {
	g.mu.Lock()
	delete(g.taken, id)
	g.mu.Unlock()
}
