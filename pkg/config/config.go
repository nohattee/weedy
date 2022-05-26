package config

import "os"

type Config struct {
	Name        string
	Environment string
	HTTP        HTTPConfig
	DB          DBConfig
}

type HTTPConfig struct {
	Host string
	Port string
}

type DBConfig struct {
	Name       string
	Host       string
	Port       string
	Username   string
	Password   string
	Connection string
}

func LoadConfigFromEnv() *Config {
	return &Config{
		Name:        os.Getenv("APP_NAME"),
		Environment: os.Getenv("APP_ENV"),
		HTTP: HTTPConfig{
			Host: os.Getenv("HTTP_HOST"),
			Port: os.Getenv("HTTP_PORT"),
		},
		DB: DBConfig{
			Name:       os.Getenv("DB_DATABASE"),
			Host:       os.Getenv("DB_HOST"),
			Port:       os.Getenv("DB_PORT"),
			Username:   os.Getenv("DB_USERNAME"),
			Password:   os.Getenv("DB_PASSWORD"),
			Connection: os.Getenv("DB_CONNECTION"),
		},
	}
}
