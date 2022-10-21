package application

import "github.com/labstack/echo/v4"

type CustomContext struct {
	echo.Context
	*Application
}

func NewCustomContext(app *Application, c echo.Context) *CustomContext {
	return &CustomContext{
		c,
		app,
	}
}

func (c *CustomContext) Echo() *echo.Echo {
	return c.Application.Echo
}
