package goui

// Children is either nil or an instance of one of these types:
//
//		*Node
//	 []*Node
//		string
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
	panic("unreachable code")
}
