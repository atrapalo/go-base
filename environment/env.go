package environment

import "errors"

const Prod = "prod"
const Staging = "staging"
const Local = "local"

var validEnvs = map[string]bool{
	Prod:    true,
	Staging: true,
	Local:   true,
}

type Env struct {
	val string
}

func NewEnv(v string) (Env, error) {
	if !isValidEnv(v) {
		return Env{}, errors.New("invalid environment: " + v)
	}

	return Env{val: v}, nil
}

func (e Env) Value() string {
	return e.val
}

func (e Env) IsProd() bool {
	return e.val == Prod
}

func (e Env) IsStaging() bool {
	return e.val == Staging
}

func (e Env) IsLocal() bool {
	return e.val == Local
}

func isValidEnv(v string) bool {
	_, ok := validEnvs[v]

	return ok
}
