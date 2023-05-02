package goui

import (
	"github.com/twharmon/godom"
)

func Mount(selector string, r *Node) {
	godom.Mount(selector, r.createDom())
	<-make(chan struct{})
}
