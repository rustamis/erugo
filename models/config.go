package models

type Config struct {
	BaseStoragePath string `json:"base_storage_path"`
	AppUrl          string `json:"app_url"`
	BindPort        int    `json:"bind_port"`
	JwtSecret       string `json:"jwt_secret"`
	MaxShareSize    string `json:"max_share_size"`
}
