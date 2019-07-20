package api

import (
	"net/http"

	"github.com/champly/hercules/ctxs"
)

func Api(ctx *ctxs.Context) (err error) {

	ctx.Log.Info("===========api demo==============")

	ctx.Log.Info("info")
	ctx.Log.Debug("debug")
	ctx.Log.Warn("warn")

	ctx.JSON(http.StatusOK, map[string]interface{}{
		"code": "succ",
		"msg":  "success",
	})
	return nil
}
