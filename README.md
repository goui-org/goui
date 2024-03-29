# GoUI
An experimental web framework for creating user interfaces.

## GoUIX
Install `gouix`, a cli tool for GoUI.

```
go install github.com/goui-org/gouix@latest
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
	"github.com/goui-org/goui"
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

	"github.com/goui-org/goui"
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
