package config

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var Global *ServerConfig

func InitConfig(path ...string) {
	var config string
	if len(path) == 0 {
		//从配置文件中读取出对应的配置
		configFilePrefix := "config"
		config = fmt.Sprintf("./%s.yaml", configFilePrefix)
	} else {
		config = path[0]
	}

	v := viper.New()
	//文件的路径如何设置
	v.SetConfigFile(config)
	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}
	v.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("config file changed:", e.Name)
		if err := v.Unmarshal(&Global); err != nil {
			fmt.Println(err)
		}
	})
	//这个对象如何在其他文件中使用 - 全局变量
	if err := v.Unmarshal(&Global); err != nil {
		panic(err)
	}
	fmt.Println("Global = ", Global)
	zap.S().Infof("配置信息: %v", Global)
}
