package configs

type ServerConfig struct {
	PlatName   string `json:"plat_name"`
	SystemName string `json:"system_name"`
	Addr       string `json:"addr"`
	Mode       string `json:"mode"`
}
