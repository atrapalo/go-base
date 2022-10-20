package environment_test

import (
	"github.com/stretchr/testify/assert"
	"gitlab.atrapalo.com/accommodations/go-base/environment"
	"testing"
)

func Test_env_vars_reader_is_created_with_debug_false(t *testing.T) {
	t.Setenv(environment.AppEnv, "local")
	t.Setenv(environment.Debug, "false")
	t.Setenv(environment.SilentMode, "false")
	t.Setenv(environment.WebPort, "8082")
	t.Setenv(environment.UseSSL, "true")

	reader := environment.NewReader()

	assert.Equal(t, "local", reader.AppEnv().Value())
	assert.Equal(t, false, reader.Debug())
	assert.Equal(t, false, reader.SilentMode())
	assert.Equal(t, "8082", reader.WebPort())
	assert.Equal(t, true, reader.UseSSL())
}

func Test_env_vars_reader_is_created_with_debug_true(t *testing.T) {
	t.Setenv(environment.AppEnv, "local")
	t.Setenv(environment.Debug, "true")
	t.Setenv(environment.SilentMode, "false")
	t.Setenv(environment.WebPort, "8082")
	t.Setenv(environment.UseSSL, "true")

	reader := environment.NewReader()

	assert.Equal(t, "local", reader.AppEnv().Value())
	assert.Equal(t, true, reader.Debug())
	assert.Equal(t, false, reader.SilentMode())
	assert.Equal(t, "8082", reader.WebPort())
	assert.Equal(t, true, reader.UseSSL())
}

func Test_when_no_debug_neither_app_env_neither_port_are_provided_then_is_prod_and_no_debug_and_default_port(t *testing.T) {
	t.Setenv(environment.SilentMode, "false")
	t.Setenv(environment.UseSSL, "true")

	reader := environment.NewReader()

	assert.Equal(t, "prod", reader.AppEnv().Value())
	assert.Equal(t, false, reader.Debug())
	assert.Equal(t, false, reader.SilentMode())
	assert.Equal(t, "80", reader.WebPort())
	assert.Equal(t, true, reader.UseSSL())
}

func Test_env_vars_reader_is_created_in_prod_env_if_no_app_env_is_provided(t *testing.T) {
	t.Setenv(environment.AppEnv, "")
	t.Setenv(environment.Debug, "false")
	t.Setenv(environment.SilentMode, "false")
	t.Setenv(environment.WebPort, "8082")
	t.Setenv(environment.UseSSL, "true")

	reader := environment.NewReader()

	assert.Equal(t, "prod", reader.AppEnv().Value())
	assert.Equal(t, false, reader.Debug())
	assert.Equal(t, false, reader.SilentMode())
	assert.Equal(t, "8082", reader.WebPort())
	assert.Equal(t, true, reader.UseSSL())
}

func Test_env_vars_reader_panics_to_initialize_if_invalid_is_provided(t *testing.T) {
	t.Setenv(environment.AppEnv, "piruleta")
	t.Setenv(environment.Debug, "false")
	t.Setenv(environment.SilentMode, "false")
	t.Setenv(environment.WebPort, "8082")
	t.Setenv(environment.UseSSL, "true")

	assert.Panics(t, func() {
		_ = environment.NewReader()
	})
}
