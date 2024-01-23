package goui

func Mount(selector string, component func(NoProps) *Node) {
	mount(createDom(Component(component, nil), ""), selector)
	select {}
}
