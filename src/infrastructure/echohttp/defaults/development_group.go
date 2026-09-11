package defaults

import (
	"mini-market/src/infrastructure/echohttp/defaults/handlers"

	"github.com/hibiken/asynqmon"
	"github.com/labstack/echo/v4"
)

type DevelopmentGroup struct {
	asynqmonHandler *asynqmon.HTTPHandler
}

// @inject
func NewDevelopmentGroup(asynqmonHandler *asynqmon.HTTPHandler) *DevelopmentGroup {
	return &DevelopmentGroup{asynqmonHandler: asynqmonHandler}
}

func (this *DevelopmentGroup) RegisterRoutes(group *echo.Group) {
	group.GET("/docs/*", handlers.SwaggerWrapHandler)
	group.Any("/monitoring/tasks/*", echo.WrapHandler(this.asynqmonHandler))
}
