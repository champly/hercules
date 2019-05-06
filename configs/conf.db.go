package configs

type DBConfig struct {
	Provider    string `json:"provider"`
	ConnString  string `json:"conn_string"`
	MaxOpen     int    `json:"max_open"`
	MaxIdle     int    `json:"max_idle"`
	MaxLifeTime int    `json:"max_life_time"`
}
