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
	case []any:
		nodes := make([]*Node, len(v))
		for i := range v {
			switch v := v[i].(type) {
			case *Node:
				nodes[i] = v
			default:
				nodes[i] = Text("%v", v)
			}
		}
		return nodes
	}
	return []*Node{Text("%v", r)}
}
