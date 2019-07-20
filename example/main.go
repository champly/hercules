package main

import (
	"github.com/champly/hercules/cmd/hercules"
	"github.com/champly/hercules/example/services/api"
	"github.com/champly/hercules/example/services/cron"
)

type App struct {
	*hercules.Hercules
}

func main() {
	app := &App{hercules.New()}

	app.register()

	app.Start()
}

func (a *App) register() {
	a.API("/api/demo", api.Api)

	a.Cron("auto.demo", cron.NewCron)
}
