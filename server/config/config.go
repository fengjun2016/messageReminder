package config

import (
	"github.com/fsnotify/fsnotify"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func Init() {
	viper.AddConfigPath("conf")
	viper.SetConfigName("config")
	viper.SetConfigType("yml")    //设置配置文件格式为YML

	// viper解析配置文件
	if err := viper.ReadInConfig(); err != nil {
		logrus.Errorf("Error reading config file, %s", err)
		panic(err) 
	}

	// 监控配置文件变化并热加载程序
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event){
		logrus.Infof("Config file changed:", e.Name)
	})
}