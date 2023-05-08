package goui

type Renderable interface {
	[]*Node | *Node | string
}

func renderableToNodes(r any) []*Node {
	if r == nil {
		return nil
	}
	switch v := r.(type) {
	case string:
		return []*Node{Text(v)}
	case *Node:
		return []*Node{v}
	case []*Node:
		return v
	}
	panic("unreachable code")
}
