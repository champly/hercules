package cron

import (
	"github.com/champly/hercules/ctxs"
)

type Cron struct{}

func NewCron() interface{} {
	return &Cron{}
}

func (c *Cron) Handler(ctx *ctxs.Context) (err error) {

	ctx.Log.Info("===========自动任务==============")

	ctx.Log.Info("info")
	ctx.Log.Debug("debug")
	ctx.Log.Warn("warn")

	// ctx.ToolBox.Produce("mq.test", fmt.Sprintf(`{"time":"%s"}`, time.Now().Format("2006-01-02 15:04:05")))

	return nil
}

func Receive(ctx *ctxs.Context) (err error) {

	ctx.Log.Info("==========Receive Mq Info=========")
	ctx.Log.Info(ctx.GetBody())

	return nil
}
