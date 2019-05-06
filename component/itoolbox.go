package component

type IToolBox interface {
}

type ToolBox struct {
}

func NewToolBox() *ToolBox {
	return &ToolBox{}
}
