package routes

import (
	"fmt"
	"github.com/gin-gonic/gin"
	middlewares "github.com/goasali/toolkit/http/middleware"
	"github.com/goasali/toolkit/http/validations"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Interface interface {
	Listen(*RouteModuleParams)
}

type RouteModule struct {
	Interface
}

func NewRouteModule() *RouteModule {
	return &RouteModule{}
}

type RouteModuleParams struct {
	Router *gin.RouterGroup
}

type Route struct {
	*gin.Engine
	config *Config
}

func SetupRouter(configFunctions ...ConfigFunc) *Route {
	appConfig := getConfig(configFunctions...)

	if appConfig.mode != "" {
		gin.SetMode(appConfig.mode)
	}
	router := gin.Default()

	router.Use(middlewares.Logging())
	//router.Use(middlewares.Recovery())

	r := &Route{router, &appConfig}
	if db := appConfig.db; db != nil {
		r.loadValidations(db)
	}

	return r
}

func (r *Route) loadValidations(db *gorm.DB) {
	if err := validations.AddDatabase(db); err != nil {
		log.Fatalf("error during load database validation: %v", err)
	}
}

func (r *Route) AddApiRoutes(routes ...Interface) {
	r.AddRoutes("/api", routes...)
}

func (r *Route) AddRoutes(routePath string, routes ...Interface) {
	for _, route := range routes {
		grp := r.Group(routePath)
		route.Listen(&RouteModuleParams{grp})
	}
}

func (r *Route) Listen() error {
	addr := fmt.Sprintf("%s:%d", r.config.host, r.config.port)
	if err := r.Run(addr); err != nil {
		return err
	}

	return nil
}
