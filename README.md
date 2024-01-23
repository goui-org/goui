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
	goui.Mount("#root", app.App)
}
```


```go
// app/app.go
package app

import (
	"strconv"

	"github.com/twharmon/goui"
)

func App(goui.NoProps) *goui.Node {
	count, setCount := goui.UseState(0)

	goui.UseEffect(func() goui.EffectTeardown {
		goui.Log("count is %d", count)
		return nil
	}, goui.Deps{count})

	handleIncrement := goui.UseCallback(func(e *goui.MouseEvent) {
		setCount(func(c int) int { return c + 1 })
	}, goui.Deps{})

	return goui.Element("div", &goui.Attributes{
		Class: "app",
		Slot: []*goui.Node{
			goui.Element("button", &goui.Attributes{
				Class:   "app-btn",
				Slot:    "increment",
				OnClick: handleIncrement,
			}),
			goui.Element("p", &goui.Attributes{
				Slot: "count: " + strconv.Itoa(count),
			}),
		},
	})
}
```
