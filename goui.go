package goui

import (
	"github.com/twharmon/godom"
)

func Mount(selector string, n *Node) {
	godom.Mount(selector, n.createDom())
	<-make(chan struct{})
}
