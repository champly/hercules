package hercules

import "github.com/champly/hercules/configs"

type option struct {
	*configs.ServerConfig
	ServiceType []string
}

type Option func(*option)

func WithPlatName(pn string) Option {
	return func(o *option) {
		o.PlatName = pn
	}
}

func WithSystemName(sn string) Option {
	return func(o *option) {
		o.SystemName = sn
	}
}

func WithAddr(addr string) Option {
	return func(o *option) {
		o.Addr = addr
	}
}

func WithMode(mode string) Option {
	return func(o *option) {
		o.Mode = mode
	}
}

func WithServerType(types ...string) Option {
	return func(o *option) {
		o.ServiceType = types
	}
}
