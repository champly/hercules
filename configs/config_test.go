package configs

import (
	"testing"

	"github.com/spf13/viper"
)

func TestSetup(t *testing.T) {
	viper.SetConfigFile("./config.yaml")
	viper.ReadInConfig()
	Setup()

	t.Logf("%+v\n", PlatInfo)
	t.Logf("%+v\n", SystemInfo)
	t.Logf("%+v\n", LoggerInfo)
	t.Logf("%+v\n", HTTPServerInfo)
	t.Logf("%+v\n", CronServerInfo)
	t.Logf("%+v\n", DBInfo)
	t.Logf("%+v\n", MQServer)
}
