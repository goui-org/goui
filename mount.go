// import { Elem, createDom } from './elem.js';

//	export let mount = (root: HTMLElement, elem: Elem) => {
//	    root.appendChild(createDom(elem));
//	};
package goui

func Mount(selector string, node *Node) {
	root := document.Call("querySelector", selector)
	root.Call("appendChild", createDom(node, ""))
	select {}
}
