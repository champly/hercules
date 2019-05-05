package configs

type Router struct {
	Name    string
	Method  string
	Args    map[string]string
	Handler interface{}
}
