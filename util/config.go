package util

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

// Config stores all configuaration of the application.
type Config struct {
	DB     Database `mapstructure:"data_source"`
	Server Server   `mapstructure:"server"`
	Token  Token    `mapstructure:"jwt"`
}

func LoadConfig(name, path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName(name)
	viper.SetConfigType("yaml")

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}

type Server struct {
	Http ServerAddress `mapstructure:"http"`
	Grpc ServerAddress `mapstructure:"grpc"`
}

func (s *ServerAddress) Addr() string {
	return fmt.Sprintf("%s:%s", s.Host, s.Port)
}

type ServerAddress struct {
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
}

type Database struct {
	Driver string `mapstructure:"driver"`
	Source string `mapstructure:"source"`
}

type Token struct {
	Secret               string        `mapstructure:"secret"`
	AccessTokenDuration  time.Duration `mapstructure:"access_token_duration"`
	RefreshTokenDuration time.Duration `mapstructure:"refresh_token_duration"`
}
