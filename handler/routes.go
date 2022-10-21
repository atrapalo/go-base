package handler

import (
	"github.com/atrapalo/go-base/engine"
)

func DefineRoutes(app *engine.Engine) {
	g := app.Group("/api")

	g.GET("/status", GetApiStatus).Name = "get-api-status"
}
