package configs

type Plat struct {
	Name string
}

type System struct {
	Name       string
	ServerType string
}

type Registry struct {
	Addr      string
	IsLocal   string
	LocalAddr string
}

type Logger struct {
	Level string
	Out   []string
	Local string
	Addr  string
}
