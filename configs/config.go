package configs

import (
	"github.com/spf13/viper"
)

type Config struct {
	Plat       plat       `json:"plat"`
	System     system     `json:"system"`
	Logger     logger     `json:"logger"`
	HttpServer httpserver `json:"httpserver"`
	CronServer cronserver `json:"cronserver"`
}

type plat struct {
	Name string `json:"name"`
}

type system struct {
	Name string `json:"name"`
	Mode string `json:"mode"`
	Type string `json:"type"`
}

type logger struct {
	Level string `json:"level"`
	Out   string `json:"out"`
}

type httpserver struct {
	Cors struct {
		Enable bool              `json:"enable"`
		Header map[string]string `json:"header"`
	} `json:"cors"`
	Address string `json:"address"`
	Status  string `json:"status"`
}

type cronserver struct {
	Status   string `json:"status"`
	TaskList []struct {
		Name string `json:"name"`
		Time string `json:"time"`
	} `json:"tasklist"`
}

var (
	PlatInfo       = &plat{}
	SystemInfo     = &system{}
	LoggerInfo     = &logger{}
	HttpServerInfo = &httpserver{}
	CronServerInfo = &cronserver{}
)

func setDefault() {
	viper.SetDefault("plat", plat{Name: "hercules-plat"})
	viper.SetDefault("system", system{Name: "hercules-system", Mode: "debug"})
	viper.SetDefault("httpserver", httpserver{Address: ":8080"})
	viper.SetDefault("logger", logger{Level: "all", Out: "stdio"})
}

func Setup() {
	setDefault()

	config := Config{}
	if err := viper.Unmarshal(&config); err != nil {
		panic("read config fail:" + err.Error())
	}

	PlatInfo = &config.Plat
	SystemInfo = &config.System
	LoggerInfo = &config.Logger
	HttpServerInfo = &config.HttpServer
	CronServerInfo = &config.CronServer
}
