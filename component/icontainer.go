package component

type IContainer interface {
}

type Container struct {
}

func NewContainer() *Container {
	return &Container{}
}
