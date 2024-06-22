package config

import (
	"log"
	"time"

	"github.com/spf13/viper"
)

var Conf *Config

type Config struct {
	DBDriver            string        `mapstructure:"DB_DRIVER"`
	DBSource            string        `mapstructure:"DB_SOURCE"`
	HTTPServerAddress   string        `mapstructure:"HTTP_SERVER_ADDRESS"`
	JWTSecretKey        string        `mapstructure:"JWT_SECRETKEY"`
	AccessTokenDuration time.Duration `mapstructure:"ACCESS_TOKEN_EXPIRY"`
}

func loadConfig() (*Config, error) {
	viper.AddConfigPath(".")
	viper.SetConfigFile(".env")

	var config Config
	err := viper.ReadInConfig()
	if err != nil {
		log.Println("error read in config")
		return &config, err
	}

	err = viper.Unmarshal(&config)
	return &config, err
}

func InitConfig() {
	var err error
	Conf, err = loadConfig()
	if err != nil {
		log.Fatal("cannot load config: ", err)
	}
}
