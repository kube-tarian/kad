package config

import (
	"github.com/spf13/viper"
	"os"
	"sync"
)

const (
	DefaultConfigPath = "/etc/server/"
)

var (
	configObj *viper.Viper
	once      sync.Once
)

func bindEnvs(cfg *viper.Viper) {
	cfg.SetDefault("server.host", "0.0.0.0")
	cfg.SetDefault("server.port", 9092)
	cfg.SetDefault("server.db", "astra")
	_ = cfg.BindEnv("server.host", "HOST")
	_ = cfg.BindEnv("server.port", "PORT")
}

func New() (*viper.Viper, error) {
	var err error
	once.Do(func() {
		cfg := viper.New()
		bindEnvs(cfg)
		cfg.SetConfigType("yaml")
		cfg.SetConfigName("config")
		cfg.AddConfigPath(getConfigPath())
		err = cfg.ReadInConfig()
		configObj = cfg
	})

	return configObj, err
}

func GetConfig() *viper.Viper {
	return configObj
}

func getConfigPath() string {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		return DefaultConfigPath
	}

	return configPath
}
