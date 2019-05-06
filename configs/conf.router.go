package configs

type Router struct {
	Name    string
	Method  string
	Cron    string
	Args    map[string]string
	Handler ExecHandler
}
