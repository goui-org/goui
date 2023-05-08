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
	"github.com/twharmon/godom"
	"github.com/twharmon/goui"
)

func App(_ any) *goui.Node {
	cnt, setCnt := goui.UseState(0)

	goui.UseEffect(func() goui.EffectTeardown {
		godom.Console.Log("count changed to %d", cnt)
		return nil
	}, cnt)

	return goui.Element("div", goui.Attributes{
		Children: []*goui.Node{
			goui.Element("button", goui.Attributes{
				Children: goui.Text("increment").Slice(),
				OnClick: func(e *godom.MouseEvent) {
					setCnt(func(c int) int { return c + 1 })
				},
			}),
			goui.Element("p", goui.Attributes{
				Children: goui.Text("cnt: %d", cnt).Slice(),
			}),
		},
	})
}
```go
