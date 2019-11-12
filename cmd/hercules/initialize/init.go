package initialize

import (
	"log"

	"github.com/google/gops/agent"
)

func Initialize() {
	if err := agent.Listen(agent.Options{
		Addr: ":16166",
	}); err != nil {
		log.Fatal(err)
	}
	return
}
