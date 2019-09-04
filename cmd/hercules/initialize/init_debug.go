//+build debug

package initialize

import (
	"log"
	"net/http"
	_ "net/http/pprof"

	"github.com/fatih/color"
)

func Initialize() {
	color.Red("开启pprof, 使用:\n\tgo tool pprof -http=:8081 http://localhost:6066/debug/pprof/heap\n查看相信信息")
	go func() {
		// terminal: $ go tool pprof -http=:8081 http://localhost:6066/debug/pprof/{heap,allocs,block,cmdline,goroutine,mutex,profile,threadcreate,trace}
		// web:
		// 1、http://localhost:8081/ui
		// 2、http://localhost:6066/debug/charts
		// 3、http://localhost:6066/debug/pprof
		log.Println(http.ListenAndServe("0.0.0.0:6066", nil))
	}()

	return
}
