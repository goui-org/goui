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
	goui.Mount("#root", goui.Component(app.App, goui.NoProps))
}
```


```go
// app/app.go
package app

import (
	"github.com/twharmon/godom"
	"github.com/twharmon/goui"
)

func App(_ any) *goui.Node {
	count, setCount := goui.UseState(0)

	goui.UseEffect(func() goui.EffectTeardown {
		godom.Console.Log("count is %d", count)
		return nil
	}, count)

	return goui.Element("div", goui.Attributes{
		Children: []*goui.Node{
			goui.Element("button", goui.Attributes{
				Children: "increment",
				OnClick: func(e *godom.MouseEvent) {
					setCount(func(c int) int { return c + 1 })
				},
			}),
			goui.Element("p", goui.Attributes{
				Children: fmt.Sprintf("count: %d", count),
			}),
		},
	})
}
```
