package hercules

import (
	"github.com/champly/hercules/configs"
)

type option struct {
	ServiceType []string
}

type Option func(*option)

func WithPlatName(pn string) Option {
	return func(o *option) {
		configs.PlatInfo.Name = pn
	}
}

func WithSystemName(sn string) Option {
	return func(o *option) {
		configs.SystemInfo.Name = sn
	}
}

func WithAddr(addr string) Option {
	return func(o *option) {
		configs.HttpServerInfo.Address = addr
	}
}

func WithMode(mode string) Option {
	return func(o *option) {
		configs.SystemInfo.Mode = mode
	}
}

func WithServerType(types ...string) Option {
	return func(o *option) {
		o.ServiceType = types
	}
}
