package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

// LoadFile reads a config file to get environment variables,
// then loads them into a struct with the appropriate tags.
// If no path to the config directory is provided an error is returned.
func LoadFile(configPath, filename string, envStruct any) error {
	if configPath == "" {
		return ErrMissingConfigPath
	}

	if filename == "" {
		return ErrMissingFilename
	}

	file := strings.Split(filename, ".")
	if len(file) != 2 {
		return fmt.Errorf("%w: a filename of this format 'name.ext' is required", ErrInvalidFilename)
	}

	name, ext := file[0], file[1]
	if name == "" {
		return fmt.Errorf("%w: %s", ErrInvalidFilename, filename)
	}
	if ext == "" {
		// Should I write validation for extension types: .json, .toml, .yaml, etc. ?
		return fmt.Errorf("%w: a valid file extension is required", ErrInvalidFilename)
	}

	viper.AddConfigPath(configPath)
	viper.SetConfigName(name)
	viper.SetConfigType(ext)

	if err := viper.ReadInConfig(); err != nil {
		// what other errors am I missing?
		if errors.As(err, new(viper.ConfigFileNotFoundError)) {
			return fmt.Errorf("%s: no such config file found in directory %s", filename, configPath)
		}

		if errors.As(err, new(viper.ConfigParseError)) {
			// check if this file exists
			if _, err := os.Stat(filepath.Join(configPath, filename)); os.IsNotExist(err) {
				return fmt.Errorf("%s: no such config file found in directory %s", filename, configPath)
			}

			// I'm considering a custom error message, but the ones that
			// viper provides are already really descriptive and useful.
			return err
		}

		return err
	}

	if err := viper.Unmarshal(&envStruct); err != nil {
		return err
	}

	return nil
}
