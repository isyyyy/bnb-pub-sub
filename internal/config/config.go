package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type Root struct {
	Server   ServerConfig   `mapstructure:"server"`
	Mongo    MongoConfig    `mapstructure:"mongo"`
	RabbitMQ RabbitMQConfig `mapstructure:"rabbitmq"`
	Bnb      BnbConfig      `mapstructure:"bnb"`
}

type ServerConfig struct {
	Name string `mapstructure:"name"`
	Port string `mapstructure:"port"`
}

type MongoConfig struct {
	URL      string `mapstructure:"url"`
	Port     string `mapstructure:"port"`
	Database string `mapstructure:"database"`
}

type RabbitMQConfig struct {
	URL          string `mapstructure:"url"`
	ExchangeName string `mapstructure:"exchange_name"`
	ExchangeKind string `mapstructure:"exchange_kind"`
	QueueName    string `mapstructure:"queue_name"`
}

type BnbConfig struct {
	URL        string   `mapstructure:"url"`
	ListSymbol []string `mapstructure:"listsymbol"`
	Kline      []string `mapstructure:"kline"`
}

func InitConfig() Root {
	//Set the file name of the configuration file
	viper.SetConfigName("config")

	//Set the path to look for the configuration file
	viper.AddConfigPath(".")

	//Enable VIPER to read Environment Variables
	viper.AutomaticEnv()

	viper.SetConfigType("yml")
	var root Root

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading root file, %s ", err)
	}

	////Set undefined variables
	//viper.SetDefault("test","test-root")

	err := viper.Unmarshal(&root)
	if err != nil {
		fmt.Printf("Unable to decode into struct, %v ", err)
	}
	return root
}
