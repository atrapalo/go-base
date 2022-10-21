package application

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/newrelic/go-agent/v3/integrations/nrecho-v4"
	"github.com/newrelic/go-agent/v3/newrelic"
	"os"
	"path/filepath"
)

type Application struct {
	*echo.Echo
	debug      bool
	port       string
	silentMode bool
	useSSL     bool
}

func New(debug bool, port string, silentMode bool, useSSL bool, newRelicApp *newrelic.Application) *Application {
	app := &Application{
		echo.New(),
		debug,
		port,
		silentMode,
		useSSL,
	}

	app.Debug = debug
	app.HideBanner = true

	if newRelicApp != nil {
		skipperConfig := nrecho.WithSkipper(requestedPathMustNotBeMonitored)
		app.Use(nrecho.Middleware(newRelicApp, skipperConfig))
	}

	app.Use(customContextMiddleware(app))
	app.Use(middleware.RequestID())
	app.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
	}))
	app.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Skipper: func(c echo.Context) bool {
			return mustDisableAccessLog(c, silentMode)
		},
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))
	app.Use(middleware.Recover())

	app.Static("/", app.getDir()+"/swagger-ui").Name = "swagger-static"
	app.File("/swagger", app.getDir()+"/swagger.yaml").Name = "swagger-file"

	return app
}

func (a *Application) Start() {
	if a.useSSL {
		a.startSSL()
	}

	a.Logger.Fatal(a.Echo.Start(":" + a.port))
}

func (a *Application) startSSL() {
	go func() {
		a.Logger.Fatal(a.Echo.StartTLS(":443", a.getDir()+"/ssl/cert.pem", a.getDir()+"/ssl/key.pem"))
	}()
}

func (a *Application) getDir() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		a.Logger.Error(err)
		dir = ""
	}

	return dir
}

func customContextMiddleware(app *Application) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := NewCustomContext(app, c)
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
