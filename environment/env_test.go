package environment_test

import (
	"errors"
	"github.com/atrapalo/go-base/environment"
	"github.com/stretchr/testify/assert"
	"testing"
)

type data struct {
	Name         string
	Request      Request
	Expectations dataExpectations
}

type Request = string

type dataExpectations struct {
	Environment string
	IsProd      bool
	IsStaging   bool
	IsLocal     bool
	Error       error
}

func dataProvider() []data {
	return []data{
		{
			Name:    "test-environment-local",
			Request: "local",
			Expectations: dataExpectations{
				Environment: "local",
				IsProd:      false,
				IsStaging:   false,
				IsLocal:     true,
				Error:       nil,
			},
		},
		{
			Name:    "test-environment-staging",
			Request: "staging",
			Expectations: dataExpectations{
				Environment: "staging",
				IsProd:      false,
				IsStaging:   true,
				IsLocal:     false,
				Error:       nil,
			},
		},
		{
			Name:    "test-environment-prod",
			Request: "prod",
			Expectations: dataExpectations{
				Environment: "prod",
				IsProd:      true,
				IsStaging:   false,
				IsLocal:     false,
				Error:       nil,
			},
		},
		{
			Name:    "test-environment-invalid",
			Request: "piruleta 🦄",
			Expectations: dataExpectations{
				Environment: "",
				IsProd:      false,
				IsStaging:   false,
				IsLocal:     false,
				Error:       errors.New("invalid environment: piruleta 🦄"),
			},
		},
	}
}

func Test_environment_is_created(t *testing.T) {
	for _, d := range dataProvider() {
		t.Run(d.Name, func(t *testing.T) {
			env, err := environment.NewEnv(d.Request)
			assert.Equal(t, d.Expectations.Environment, env.Value())
			assert.Equal(t, d.Expectations.IsLocal, env.IsLocal())
			assert.Equal(t, d.Expectations.IsStaging, env.IsStaging())
			assert.Equal(t, d.Expectations.IsProd, env.IsProd())
			assert.Equal(t, d.Expectations.Error, err)
			if err != nil {
				assert.Empty(t, env)
				assert.Equal(t, d.Expectations.Error.Error(), err.Error())
			}
		})
	}
}
