package configs

// HttpMethod
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

// ServiceType support service type
const (
	ServerTypeAPI  = "api"
	ServerTypeCron = "cron"
	ServerTypeMQ   = "mq"
)

// ServiceMode service run mode
const (
	ServerModeDebug   = "debug"
	ServerModeRelease = "release"
)
