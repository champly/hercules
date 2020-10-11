package signal

import (
	"os"
	"os/signal"
	"syscall"

	"k8s.io/klog/v2"
)

func SetupSignalHandler() <-chan struct{} {

	stop := make(chan struct{})

	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM)
	go func() {
		<-c
		close(stop)
		<-c
		klog.Warning("force close!")
		os.Exit(-1)
	}()

	return stop
}
