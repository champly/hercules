//+build !debug

package initialize

import (
	"log"

	"github.com/google/gops/agent"
)

// Setpprof set pprof service
func Setpprof() {
	if err := agent.Listen(agent.Options{
		Addr: ":16166",
	}); err != nil {
		log.Fatal(err)
	}
	return
}
