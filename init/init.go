package init

import (
	"fmt"

	_ "github.com/champly/hercules/servers/cron"
	_ "github.com/champly/hercules/servers/http"
	"github.com/spf13/viper"
)

func init() {
	fmt.Println("init execute")
	initConfig()
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.AddConfigPath("$HOME/config")
	viper.AddConfigPath("./")
	viper.AddConfigPath("./config")
	viper.SetConfigName("config")

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		// panic("load config file error:" + err.Error())
		return
	}

	fmt.Println("Using config file:", viper.ConfigFileUsed())
	return
}
