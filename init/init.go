package init

import (
	"fmt"

	"github.com/champly/hercules/configs"
	_ "github.com/champly/hercules/servers/cron"
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

	fmt.Println("Using config file:", viper.ConfigFileUsed())
	return
}
