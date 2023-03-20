package config

import (
	"errors"
	"fmt"
)

var ErrNoEnvs = errors.New("no environment variables provided")

type EnvConfigErr struct {
	missingEnvs []string
}

func (e *EnvConfigErr) Error() string {
	if !e.hasErrors() {
		return ""
	}

	var errs string
	if len(e.missingEnvs) != 0 {
		errs += fmt.Sprintf("missing %v", e.missingEnvs)
	}

	errs += " environment variable(s)"
	return errs
}

func (m *EnvConfigErr) addMissingEnv(missing string) {
	m.missingEnvs = append(m.missingEnvs, missing)
}

func (m *EnvConfigErr) hasErrors() bool {
	return len(m.missingEnvs) != 0
}
