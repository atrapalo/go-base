package engine

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/newrelic/go-agent/v3/integrations/nrecho-v4"
	"github.com/newrelic/go-agent/v3/newrelic"
	"os"
	"path/filepath"
)

type Engine struct {
	*echo.Echo
	debug      bool
	port       string
	silentMode bool
	useSSL     bool
}

func New(debug bool, port string, silentMode bool, useSSL bool, newRelicApp *newrelic.Application) *Engine {
	engine := &Engine{
		echo.New(),
		debug,
		port,
		silentMode,
		useSSL,
	}

	engine.Debug = debug
	engine.HideBanner = true

	if newRelicApp != nil {
		skipperConfig := nrecho.WithSkipper(requestedPathMustNotBeMonitored)
		engine.Use(nrecho.Middleware(newRelicApp, skipperConfig))
	}

	engine.Use(customContextMiddleware(engine))
	engine.Use(middleware.RequestID())
	engine.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
	}))
	engine.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Skipper: func(c echo.Context) bool {
			return mustDisableAccessLog(c, silentMode)
		},
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))
	engine.Use(middleware.Recover())

	engine.Static("/", engine.getDir()+"/swagger-ui").Name = "swagger-static"
	engine.File("/swagger", engine.getDir()+"/swagger.yaml").Name = "swagger-file"

	return engine
}

func (a *Engine) Start() {
	if a.useSSL {
		a.startSSL()
	}

	a.Logger.Fatal(a.Echo.Start(":" + a.port))
}

func (a *Engine) startSSL() {
	go func() {
		a.Logger.Fatal(a.Echo.StartTLS(":443", a.getDir()+"/ssl/cert.pem", a.getDir()+"/ssl/key.pem"))
	}()
}

func (a *Engine) getDir() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		a.Logger.Error(err)
		dir = ""
	}

	return dir
}

func customContextMiddleware(engine *Engine) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := NewCustomContext(engine, c)
			return next(cc)
		}
	}
}

func mustDisableAccessLog(c echo.Context, isSilentMode bool) bool {
	return isSilentMode || requestedPathMustNotBeMonitored(c)
}

func requestedPathMustNotBeMonitored(c echo.Context) bool {
	return c.Path() == "/" ||
		c.Path() == "/*" ||
		c.Path() == "/index.html" ||
		c.Path() == "/swagger" ||
		c.Path() == "/api/status"
}
