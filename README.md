# GoUI
Experimental web ui library for creating user interfaces.

## Usage
```
// main.go
package main

import (
	"github.com/twharmon/goui"
	"main/app"
)

func main() {
	goui.Mount("#root", goui.Component(app.App, goui.NoProps))
}
```go

```
// app/app.go
package app

import (
	"github.com/twharmon/goui"
	"main/app"
)

func App(_ any) *goui.Node {
	
}
```go
