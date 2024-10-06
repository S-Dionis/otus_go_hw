package main

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
)

type Config struct {
	Logger LoggerConf
	Server ServerConf
	DbType string
}

type LoggerConf struct {
	Level string `mapstructure:"level"`
}
type ServerConf struct {
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
}

type DbType struct {
	Type string `mapstructure:"db"`
}

func NewConfig(path string) Config {
	viper.SetConfigFile(path)

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Printf("Error reading config file, %s", err)
		os.Exit(1)
	}

	var loggerConf LoggerConf
	var serverConf ServerConf
	err = viper.Sub("logger").Unmarshal(&loggerConf)
	if err != nil {
		fmt.Printf("Error unmarshalling config file, %s", err)
		os.Exit(1)
	}
	err = viper.Sub("server").Unmarshal(&serverConf)
	if err != nil {
		fmt.Printf("Error unmarshalling config file, %s", err)
		os.Exit(1)
	}
	return Config{Logger: loggerConf, Server: serverConf}
}
