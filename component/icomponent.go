package component

type IComponent interface {
	GetRouter() []Router
	Close() error
}

type Router struct {
	Name   string
	Handle interface{}
}
