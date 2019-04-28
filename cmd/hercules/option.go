package hercules

type option struct {
	Addr string
	Mode string
}

type Option func(*option)

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
