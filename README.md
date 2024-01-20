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
