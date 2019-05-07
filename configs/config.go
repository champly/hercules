package configs

import (
	"github.com/spf13/viper"
)

type ServerConfig struct {
	PlatName   string `json:"plat_name"`
	SystemName string `json:"system_name"`
	Addr       string `json:"addr"`
	Mode       string `json:"mode"`
}

type plat struct {
	Name string `json:"name"`
}

var PlatInfo = &plat{}

type system struct {
	Name string `json:"name"`
	Mode string `json:"mode"`
}

var SystemInfo = &system{}

type httpserver struct {
	CorsConfig httpcors `json:"cors_config"`
	Address    string   `json:"address"`
}

type httpcors struct {
	NeedCors bool                `json:"need_cors"`
	Header   []map[string]string `json:"header"`
}

var HttpServerInfo = &httpserver{}

type log struct {
	LogLevel string `json:"log_level"`
	Out      string `json:"out"`
}

var LogInfo = &log{}

func Setup() {
	viper.Get("plat")
}
