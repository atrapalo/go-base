package engine

import "github.com/labstack/echo/v4"

type CustomContext struct {
	echo.Context
	*Engine
}

func NewCustomContext(engine *Engine, c echo.Context) *CustomContext {
	return &CustomContext{
		c,
		engine,
	}
}

func (c *CustomContext) Echo() *echo.Echo {
	return c.Engine.Echo
}
