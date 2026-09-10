package handlers

import (
	"mini-market/src/infrastructure/docs"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func SwaggerWrapHandler(c echo.Context) error {
	host := c.Request().Header.Get("X-Forwarded-Host")
	if host == "" {
		host = c.Request().Host
	}
	docs.SwaggerInfo.Host = host
	docs.SwaggerInfo.BasePath = "/api/"

	handler := echoSwagger.EchoWrapHandler(echoSwagger.PersistAuthorization(true))
	return handler(c)
}
