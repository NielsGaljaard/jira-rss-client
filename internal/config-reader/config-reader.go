package config_reader

import (
	"fmt"

	"github.com/NielsGaljaard/jira-rss-client/internal/domain"
	"github.com/spf13/viper"
)

type Config struct {
	TicketsPerJQL    uint32           `mapstructure:"TicketsPerJQL"`
	TemplateLocation string           `mapstructure:"TemplateLocation"`
	AppPassword      string           `mapstructure:"AppPassword"`
	AppUrl           string           `mapstructure:"AppURL"`
	AppUser          string           `mapstructure:"AppUser"`
	LogLevel         string           `mapstructure:"LogLevel"`
	Channel          []domain.Channel `mapstructure:"Channels"`
}

func LoadConfig(path string) (*Config, error) {
	var conf Config
	viper.AddConfigPath(path)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.SetDefault("TicketsPerJQL", 50)
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("can't read config")
		return nil, err
	}
	err = viper.Unmarshal(&conf)
	if err != nil {
		fmt.Println("can't unmarshall config")
		return nil, err
	}
	return &conf, nil
}
