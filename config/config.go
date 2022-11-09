package config

import "github.com/spf13/viper"

// Config
type Config struct {
	Debug    bool     `mapstructure:"debug"`
	Server   Server   `mapstructure:"server"`
	Context  Context  `mapstructure:"context"`
	Database Database `mapstructure:"database"`
}

type Server struct {
	Address string `mapstructure:"address"`
}

type Context struct {
	Timeout int `mapstructure:"timeout"`
}

type Database struct {
	DBHost string `mapstructure:"host"`
	DBPort string `mapstructure:"port"`
	DBUser string `mapstructure:"user"`
	DBPass string `mapstructure:"pass"`
	DBName string `mapstructure:"name"`
}

var vp *viper.Viper

// Load Config from JSON into stucture ...
func LoadConfig() (Config, error) {
	vp = viper.New()

	var config Config

	vp.SetConfigName("config")
	vp.SetConfigType("json")
	vp.AddConfigPath(".")
	vp.AddConfigPath("./config")

	err := vp.ReadInConfig()
	if err != nil {
		return Config{}, err
	}

	err = vp.Unmarshal(&config)
	if err != nil {
		return Config{}, err
	}
	return config, err
}
