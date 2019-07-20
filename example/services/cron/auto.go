package cron

import "github.com/champly/hercules/ctxs"

type Cron struct{}

func NewCron() interface{} {
	return &Cron{}
}

func (c *Cron) Handler(ctx *ctxs.Context) (err error) {

	ctx.Log.Info("===========自动任务==============")

	ctx.Log.Info("info")
	ctx.Log.Debug("debug")
	ctx.Log.Warn("warn")

	return nil
}
