package goui

import (
	"github.com/twharmon/godom"
)

func Mount(selector string, n *Node) {
	n.createDom()
	godom.Mount(selector, n.dom)
	<-make(chan struct{})
}
