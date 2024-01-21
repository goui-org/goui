// import { Elem, createDom } from './elem.js';

package goui

func Mount(selector string, node *Node) {
	mount(createDom(node, ""), selector)
	select {}
}
