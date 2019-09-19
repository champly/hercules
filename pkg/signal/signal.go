package signal

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/fatih/color"
)

func SetupSignalHandler() <-chan struct{} {

	stop := make(chan struct{})

	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM)
	go func() {
		<-c
		close(stop)
		<-c
		color.HiRed("强制关闭")
		os.Exit(-1)
	}()

	return stop
}
