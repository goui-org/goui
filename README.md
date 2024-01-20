# GoUI
Experimental web ui library for creating user interfaces.

## GoUIX
Create an app with the goui cli tool GoUIX

```
go install github.com/twharmon/gouix@latest
```

Create a new app
```
gouix create my-app
```

Start the development server
```
gouix serve
```

Create a production build
```
gouix build
```

## Usage
```go
// main.go
package main

import (
	"github.com/twharmon/goui"
	"main/app"
)

func main() {
	goui.Mount("#root", goui.Component(app.App, nil))
}
```


```go
// app/app.go
package app

import (
	"fmt"

	"github.com/twharmon/goui"
)

func App(goui.NoProps) *goui.Node {
	count, setCount := goui.UseState(0)

	goui.UseEffect(func() goui.EffectTeardown {
		goui.Console.Log("count is %d", count)
		return nil
	}, goui.Deps{count})

	handleIncrement := goui.UseCallback(func(e *goui.MouseEvent) {
		setCount(func(c int) int { return c + 1 })
	}, goui.Deps{})

	return goui.Element("div", &goui.Attributes{
		Class: "app",
		Children: goui.Children{
			goui.Element("button", &goui.Attributes{
				Class:    "app-btn",
				Children: goui.Children{goui.Text("increment")},
				OnClick: handleIncrement,
			}),
			goui.Element("p", &goui.Attributes{
				Children: goui.Children{goui.Text(fmt.Sprintf("count: %d", count))},
			}),
		},
	})
}
```


JS funcs:
```


let wasm;
const elements = {};
window.elements = elements;
let elementId = 1;
const createElement = tag => {
    const id = elementId++;
    elements[id] = document.createElement(tag);
    return id;
};
const createTextNode = tag => {
    const id = elementId++;
    elements[id] = document.createTextNode(tag);
    return id;
};

const go=new Go()
go.importObject.env = {
    createElementJS: (addr, length) => {
        const memory = wasm.exports.memory
        const bytes = memory.buffer.slice(addr, addr + length)
        const text = String.fromCharCode.apply(null, new Int8Array(bytes))
        return createElement(text);
    },
    createTextNodeJS: (addr, length) => {
        const memory = wasm.exports.memory
        const bytes = memory.buffer.slice(addr, addr + length)
        const text = String.fromCharCode.apply(null, new Int8Array(bytes))
        return createTextNode(text);
    },
    appendChild: (parent, child) => {
        elements[parent].appendChild(elements[child]);
    },
    setClass: (node, addr, length) => {
        const memory = wasm.exports.memory
        const bytes = memory.buffer.slice(addr, addr + length)
        const text = String.fromCharCode.apply(null, new Int8Array(bytes))
        elements[node].className = text;
    },
    setID: (node, addr, length) => {
        const memory = wasm.exports.memory
        const bytes = memory.buffer.slice(addr, addr + length)
        const text = String.fromCharCode.apply(null, new Int8Array(bytes))
        elements[node].id = text;
    },
    set: (node, addr, length, addr2, length2) => {
        const memory = wasm.exports.memory
        const bytes = memory.buffer.slice(addr, addr + length)
        const text = String.fromCharCode.apply(null, new Int8Array(bytes))

        const bytes2 = memory.buffer.slice(addr2, addr2 + length2)
        const text2 = String.fromCharCode.apply(null, new Int8Array(bytes2))

        elements[node][text] = text2
    },
    clearChildren: (node) => {
        elements[node].textContent = '';
    },
    mount: (node) => {
        document.querySelector('#root').appendChild(elements[node])
    },
}
const fetched=fetch("main.wasm");"instantiateStreaming"in WebAssembly?WebAssembly.instantiateStreaming(fetched,go.importObject).then(e=>{
    wasm = e.instance;
    go.run(e.instance)
}):fetched.then(e=>e.arrayBuffer()).then(e=>WebAssembly.instantiate(e,go.importObject).then(e=>go.run(e.instance)))</script>

```