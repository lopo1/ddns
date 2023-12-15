package config

type Aliyun struct {
	Id     string `mapstructure:"id" json:"id"`
	Secret string `mapstructure:"secret" json:"secret"`
	Demain string `mapstructure:"demain" json:"demain"`
}
