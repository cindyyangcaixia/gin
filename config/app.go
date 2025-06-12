package config

import (
	"log"

	"github.com/spf13/viper"
)

type AppConfig struct {
	MongoConfig  MongoDBConfig `mapstructure:",squash"`
	LoggerConfig LoggerConfig  `mapstructure:",squash"`
}

func LoadEnvConfig() (AppConfig, error) {
	viper.SetConfigName("env")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("read config failed: %v", err)
	}

	var cfg AppConfig
	if err := viper.Unmarshal(&cfg); err != nil {
		return AppConfig{}, err
	}

	return cfg, nil
}
