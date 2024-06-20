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

func LoadConfig(filename string) (*Config, error) {
	log.Println(filename)
	viper.SetConfigFile(filename)

	var config Config
	err := viper.ReadInConfig()
	if err != nil {
		log.Println("error read in config")
		return &config, err
	}

	err = viper.Unmarshal(&config)
	return &config, err
}
