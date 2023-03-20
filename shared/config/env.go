// Package config contains custom errors for managing environment variables
package config

import (
	"fmt"
)

type EnvConfigErr struct {
	missingEnvs []string
	invalidEnvs []string
}

func (e *EnvConfigErr) Error() string {
	if !e.HasErrors() {
		return ""
	}

	var errs string
	if len(e.missingEnvs) != 0 {
		errs += fmt.Sprintf("missing %v", e.missingEnvs)
		if len(e.invalidEnvs) != 0 {
			errs += " and "
		}
	}
	if len(e.invalidEnvs) != 0 {
		errs += fmt.Sprintf("invalid %v", e.invalidEnvs)
	}

	errs += " environment variable(s)"
	return errs
}

func (m *EnvConfigErr) AddMissingEnv(missing string) {
	m.missingEnvs = append(m.missingEnvs, missing)
}

func (m *EnvConfigErr) AddInvalidEnv(invalid string) {
	m.invalidEnvs = append(m.invalidEnvs, invalid)
}

func (m *EnvConfigErr) HasErrors() bool {
	return (len(m.missingEnvs) + len(m.invalidEnvs)) != 0
}
