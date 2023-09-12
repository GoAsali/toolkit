package config

import (
	"fmt"
	"github.com/caarlos0/env/v8"
	"github.com/gin-gonic/gin"
	"github.com/goasali/toolkit/utils/slices"
)

type AppMode string

const (
	ReleaseMode = gin.ReleaseMode
	DebugMode   = gin.DebugMode
	TestMode    = gin.TestMode
)

type App struct {
	Name string `env:"APP_NAME"`
	Host string `env:"APP_HOST"`
	Port string `env:"APP_PORT"`
	Url  string `env:"APP_URL"`
	Mode string `env:"APP_MODE"`
}

var appConfig *App

func GetApp() (*App, error) {
	if appConfig != nil {
		return appConfig, nil
	}
	appConfig = &App{}
	if err := env.Parse(appConfig); err != nil {
		return nil, err
	}
	if appConfig.Port == "" {
		appConfig.Port = "9000"
	}
	if appConfig.Mode == "" {
		appConfig.Mode = DebugMode
	}
	if appConfig.Url == "" {
		appConfig.Url = "/"
	}

	if appConfig.Url[len(appConfig.Url)-1] != '/' {
		appConfig.Url += "/"
	}

	appModes := []string{ReleaseMode, DebugMode, TestMode}
	if !slices.Contains(appModes, appConfig.Mode) {
		return nil, fmt.Errorf("invalid app mode \"%s\"", appConfig.Mode)
	}
	return appConfig, nil
}

func (app *App) GetUrl(route string) string {
	baseRoute := "/"
	if app != nil {
		baseRoute = app.Url
	}
	return baseRoute + route
}
