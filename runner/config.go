package main

import "github.com/spf13/viper"

// LoadTOMLConfig reads a TOML file and loads its
// contents to run the microservices of this application.
// filepath is the location of the config file
// filename is the name of the toml config file without the extension
// config is the struct that stores the parsed config contents
func LoadTOMLConfig(filepath, filename string, envs any) error {
	if filepath == "" {
		filepath = "."
	}

	viper.AddConfigPath(filepath)
	viper.SetConfigName(filename)
	viper.SetConfigType("toml")

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	if err := viper.Unmarshal(&envs); err != nil {
		return err
	}

	return nil
}
