package status

import (
	"github.com/champly/hercules/ctxs"
	"k8s.io/klog/v2"
)

// GetServerStatus get service health
func GetServerStatus(ctx *ctxs.Context) error {
	// TODO: not implement
	klog.Warningln("server health not finish!")
	return nil
}
