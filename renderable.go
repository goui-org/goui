package goui

type Children any

func renderableToNodes(r Children) []*Node {
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
	return []*Node{Text("%v", r)}
}
