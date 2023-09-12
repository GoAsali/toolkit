package config

import "github.com/caarlos0/env/v8"

var adminConfig *AdminConfig

type AdminConfig struct {
	Username string `env:"ADMIN_USERNAME"`
	Password string `env:"ADMIN_PASSWORD"`
}

func GetAdmin() (*AdminConfig, error) {
	if adminConfig != nil {
		return adminConfig, nil
	}

	adminConfig = &AdminConfig{}

	if err := env.Parse(adminConfig); err != nil {
		return nil, err
	}

	return adminConfig, nil
}
