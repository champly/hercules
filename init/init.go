package init

import (
	"github.com/champly/hercules/configs"
	_ "github.com/champly/hercules/servers/cron"
	_ "github.com/champly/hercules/servers/http"
	_ "github.com/champly/hercules/servers/mq"
	"github.com/fatih/color"
	"github.com/spf13/viper"
)

func init() {
	initConfig()
	configs.Setup()
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.AddConfigPath("$HOME/config")
	viper.AddConfigPath("./")
	viper.AddConfigPath("./config")
	viper.SetConfigName("config")

	if err := viper.ReadInConfig(); err != nil {
		panic("load config file error:" + err.Error())
	}

	color.HiMagenta("Using config file:%s", viper.ConfigFileUsed())
	return
}
