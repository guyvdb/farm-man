package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	IOTAddress   string
	CLIAddress   string
	DatabaseURI  string
	DatabaseName string
}

func NewConfig(filename string) (*Config, error) {
	v, err := readConfig(filename, defaults())
	if err != nil {
		return nil, err
	}

	c := &Config{
		IOTAddress:   v.GetString("iot-address"),
		CLIAddress:   v.GetString("cli-address"),
		DatabaseURI:  v.GetString("database-uri"),
		DatabaseName: v.GetString("database-name"),
	}
	return c, nil
}

func defaults() map[string]interface{} {
	return map[string]interface{}{
		"iot-address":   "192.168.8.100:6000",
		"cli-address":   "192.168.8.100:6001",
		"database-uri":  "mongodb://localhost:27017",
		"database-name": "farm-man",
	}
}

func readConfig(filename string, defaults map[string]interface{}) (*viper.Viper, error) {
	v := viper.New()
	for key, value := range defaults {
		v.SetDefault(key, value)
	}
	v.SetConfigName(filename)
	v.AddConfigPath(".")
	v.AutomaticEnv()
	err := v.ReadInConfig()
	return v, err
}
