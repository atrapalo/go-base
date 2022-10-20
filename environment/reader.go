package environment

import (
	"os"
)

const AppEnv = "APP_ENV"
const Debug = "DEBUG"
const SilentMode = "SILENT_MODE"
const WebPort = "WEB_PORT"
const UseSSL = "USE_SSL"

const defaultPort = "80"

type Reader struct {
	appEnv     Env
	debug      bool
	silentMode bool
	webPort    string
	useSSL     bool
}

func NewReader() *Reader {
	return &Reader{
		appEnv:     readEnv(),
		debug:      readDebugMode(),
		silentMode: readSilentMode(),
		webPort:    readWebPort(),
		useSSL:     readUseSSL(),
	}
}

func (r *Reader) AppEnv() Env {
	return r.appEnv
}

func (r *Reader) Debug() bool {
	return r.debug
}

func (r *Reader) SilentMode() bool {
	return r.silentMode
}

func (r *Reader) WebPort() string {
	return r.webPort
}

func (r *Reader) UseSSL() bool {
	return r.useSSL
}

func readEnv() Env {
	if os.Getenv(AppEnv) == "" {
		env, _ := NewEnv(Prod)
		return env
	}
	env, err := NewEnv(os.Getenv(AppEnv))

	if err != nil {
		panic(err.Error())
	}

	return env
}

func readDebugMode() bool {
	return 0 < len(os.Getenv(Debug)) && os.Getenv(Debug) == "true"
}

func readSilentMode() bool {
	return 0 < len(os.Getenv(SilentMode)) && os.Getenv(SilentMode) == "true"
}

func readWebPort() string {
	if os.Getenv(WebPort) != "" {
		return os.Getenv(WebPort)
	}

	return defaultPort
}

func readUseSSL() bool {
	return 0 < len(os.Getenv(UseSSL)) && os.Getenv(UseSSL) == "true"
}
