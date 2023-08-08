package routes

import "gorm.io/gorm"

type Config struct {
	mode string
	host string
	port int
	db   *gorm.DB
}

type ConfigFunc func(config *Config)

func WithHost(host string) ConfigFunc {
	return func(config *Config) {
		config.host = host
	}
}

func WithPort(port int) ConfigFunc {
	return func(config *Config) {
		config.port = port
	}
}

func WithAppMode(appMode string) ConfigFunc {
	return func(config *Config) {
		config.mode = appMode
	}
}

func WithDatabase(db *gorm.DB) ConfigFunc {
	return func(config *Config) {
		config.db = db
	}
}

func getConfig(functions ...ConfigFunc) Config {
	config := Config{
		mode: "debug",
		host: "",
		port: 9000,
	}
	for _, configFunc := range functions {
		configFunc(&config)
	}
	return config
}
