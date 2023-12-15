package config

type ServerConfig struct {
	Name   string   `mapstructure:"name" json:"name"`
	Ip     string   `mapstructure:"ip" json:"ip"`
	Urls   []string `mapstructure:"api_url" json:"api_url"`
	Aliyun Aliyun   `mapstructure:"aliyun" json:"aliyun"`
}
