package component

type IComponent interface {
	GetRouter() []map[string]interface{}
	Close() error
}
