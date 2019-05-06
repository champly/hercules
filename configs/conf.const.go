package configs

import "github.com/champly/hercules/context"

const (
	HttpMethodGet     = "GET"
	HttpMethodPost    = "POST"
	HttpMethodPut     = "PUT"
	HttpMethodPatch   = "PATCH"
	HttpMethodHead    = "HEAD"
	HttpMethodOptions = "OPTIONS"
	HttpMethodDelete  = "DELETE"
	HttpMethodConnect = "CONNECT"
	HttpMethodTrace   = "TRACE"
	HttpMethodALL     = "GET|POST|PUT|PATCH|HEAD|OPTIONS|DELETE|CONNECT|TRACE"
)

const (
	ServerTypeAPI  = "api"
	ServerTypeCron = "cron"
)

type ExecHandler func(*context.Context) error
