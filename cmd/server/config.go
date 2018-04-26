package main

import (
	"github.com/spf13/viper"
	"github.com/txchuyen/elog.stats/config"
)

func init() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
}

func loadConfig() *config.Config {
	if err := viper.ReadInConfig(); err != nil {
		panic("config: config file not found" + err.Error())
	}
	var conf config.Config
	if err := viper.Unmarshal(&conf); err != nil {
		panic("config: unable to load" + err.Error())
	}
	return &conf
}
