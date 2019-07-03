package configs

import (
	"github.com/spf13/viper"
)

type Config struct {
	Plat       plat       `json:"plat"`
	System     system     `json:"system"`
	Logger     logger     `json:"logger"`
	DB         db         `json:"db"`
	HttpServer httpserver `json:"httpserver"`
	CronServer cronserver `json:"cronserver"`
}

type plat struct {
	Name string `json:"name"`
}

type system struct {
	Name   string `json:"name"`
	Mode   string `json:"mode"`
	Type   string `json:"type"`
	Health bool   `json:"health"`
}

type logger struct {
	Level string `json:"level"`
	Out   string `json:"out"`
}

type db struct {
	List []struct {
		Name        string `json:"name"`
		Default     bool   `json:"default"`
		Provider    string `json:"provider"`
		ConnString  string `json:"connstring"`
		MaxOpen     int    `json:"maxopen"`
		MaxIdle     int    `json:"maxidle"`
		MaxLifeTime int    `json:"maxlifetime"`
	} `json:"list"`
}

type httpserver struct {
	Cors struct {
		Enable bool              `json:"enable"`
		Header map[string]string `json:"header"`
	} `json:"cors"`
	Address string `json:"address"`
}

type cronserver struct {
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
	DBInfo         = &db{}
)

func setDefault() {
	viper.SetDefault("plat", plat{Name: "hercules-plat"})
	viper.SetDefault("system", system{Name: "hercules-system", Mode: "debug", Health: true})
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
	DBInfo = &config.DB
}
