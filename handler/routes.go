package handler

import (
	"github.com/atrapalo/go-base/application"
)

func DefineRoutes(app *application.Application) {
	g := app.Group("/api")

	g.GET("/status", GetApiStatus).Name = "get-api-status"
}
