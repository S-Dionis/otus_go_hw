package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	Logger   LoggerConf
	Server   ServerConf
	DBType   DBType
	GRPCConf GRPCConf
}

type LoggerConf struct {
	Level string `mapstructure:"level"`
}
type ServerConf struct {
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
}

type DBType struct {
	Type string `mapstructure:"db"`
}

type GRPCConf struct {
	Port string `mapstructure:"port"`
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
	var grpcConf GRPCConf
	var dbConf DBType

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
	err = viper.Sub("db").Unmarshal(&dbConf)
	if err != nil {
		fmt.Printf("Error unmarshalling config file, %s", err)
		os.Exit(1)
	}
	err = viper.Sub("grpc").Unmarshal(&grpcConf)
	if err != nil {
		fmt.Printf("Error unmarshalling config file, %s", err)
		os.Exit(1)
	}

	return Config{Logger: loggerConf, Server: serverConf, DBType: dbConf, GRPCConf: grpcConf}
}
