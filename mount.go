// import { Elem, createDom } from './elem.js';

//	export let mount = (root: HTMLElement, elem: Elem) => {
//	    root.appendChild(createDom(elem));
//	};
package goui

func Mount(selector string, elem *Elem) {
	root := document.Call("querySelector", selector)
	root.Call("appendChild", createDom(elem, ""))
	<-make(chan struct{})
}
