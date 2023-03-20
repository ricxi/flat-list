package config

import (
	"errors"
	"os"
	"strconv"
)

type envMap map[string]string

func LoadEnvs(envs ...string) (envMap, error) {
	if len(envs) == 0 {
		return nil, ErrNoEnvs
	}

	errs := EnvConfigErr{}
	em := envMap{}
	for _, envar := range envs {
		enval := os.Getenv(envar)
		if enval == "" {
			errs.addMissingEnv(enval)
		} else {
			em[envar] = enval
		}
	}

	if errs.hasErrors() {
		return nil, &errs
	}

	return em, nil
}

// Checks if the value for an environment variable's key can be
// converted into an integer. This is useful for certain networking
// methods which require a port parameter to be an integer.
func (em envMap) ValidateAsInt(envkey string) error {
	enval, ok := em[envkey]
	if !ok {
		return errors.New("invalid key: this env was not found in this map")
	}

	_, err := strconv.Atoi(enval)
	return err
}
