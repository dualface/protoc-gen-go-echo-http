package main

import (
	"example/api"
	"example/docs"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func main() {
	e := echo.New()
	e.Debug = true
	e.Use(middleware.LoggerWithConfig(middleware.DefaultLoggerConfig))
	e.Logger.SetLevel(log.DEBUG)

	group := e.Group("/v1")
	api.BindAuthService(group, &api.AuthServiceImplement{})

	docs.SwaggerInfo.BasePath = "/v1"
	e.GET("/v1/swagger/*", echoSwagger.WrapHandler)

	e.Logger.Fatal(e.Start(":12345"))
}
