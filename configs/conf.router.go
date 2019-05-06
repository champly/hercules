package configs

import "github.com/champly/hercules/ctxs"

type Router struct {
	Name    string
	Method  string
	Cron    string
	Args    map[string]string
	Handler ctxs.Handler
}
