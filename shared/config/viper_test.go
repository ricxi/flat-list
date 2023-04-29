package config

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadFile(t *testing.T) {
	type args struct {
		configPath string
		filename   string
		envStruct  any
	}
	type expected struct {
		envData   string
		errString string
	}
	tests := []struct {
		name     string
		args     args
		expected expected
		// toString converts the parsed envStruct into a string,
		// so it can be easily compared to its expected value.
		// Because the tested function may be working with various file types,
		// this allows us more freedom in comparing various file types.
		// For tests where this function will not be called because an
		// error is returned before the envStruct can be used, it will
		// be defined to return an empty string in order to avoid a
		// situation where a nil pointer derefence may occur.
		toString func(t testing.TB, envStruct any) string
	}{
		{
			name: "SuccessJSONFile",
			args: args{
				configPath: "./testdata",
				filename:   "config1.json",
				envStruct: struct {
					Port string `mapstructure:"PORT" json:"PORT"`
					Host string `mapstructure:"HOST" json:"HOST"`
				}{},
			},
			expected: expected{
				envData: `{"HOST":"127.0.0.1", "PORT":"8080"}`,
			},
			toString: func(t testing.TB, envStruct any) string {
				envData, err := json.Marshal(&envStruct)
				if err != nil {
					t.Fatal(err)
				}
				return string(envData)
			},
		},
		{
			name: "SuccessTOMLFile",
			args: args{
				configPath: "./testdata",
				filename:   "config2.toml",
				envStruct: struct {
					Host string `mapstructure:"host" json:"host"`
					Port string `mapstructure:"port" json:"port"`
				}{},
			},
			expected: expected{
				envData: `{"host":"localhost", "port":"9000"}`,
			},
			toString: func(t testing.TB, envStruct any) string {
				envData, err := json.Marshal(&envStruct)
				if err != nil {
					t.Fatal(err)
				}
				return string(envData)
			},
		},
		{
			name: "FailNoPathToConfig",
			args: args{
				configPath: "",
				filename:   "config0.json", // this is just a placeholder
				envStruct:  nil,
			},
			expected: expected{
				errString: "missing config path: a valid path to the config directory is required",
			},
			toString: func(t testing.TB, envStruct any) string {
				// this function won't be called, but an empty string
				// is returned just in case it might be called to avoid a
				// nil pointer derefence error
				return ""
			},
		},
		{
			name: "FailNoFilename",
			args: args{
				configPath: "./testdata",
				filename:   "",
				envStruct:  nil,
			},
			expected: expected{
				errString: "missing filename: a valid filename is required",
			},
			toString: func(t testing.TB, envStruct any) string {
				// see previous
				return ""
			},
		},
		{
			name: "FailInvalidFilenameEmptyName",
			args: args{
				configPath: "./testdata",
				filename:   ".toml",
				envStruct:  nil,
			},
			expected: expected{
				errString: "invalid filename: .toml",
			},
			toString: func(t testing.TB, envStruct any) string {
				// see previous
				return ""
			},
		},
		{
			name: "FailInvalidFilenameNoExtension",
			args: args{
				configPath: "./testdata",
				filename:   "config.",
				envStruct:  nil,
			},
			expected: expected{
				errString: "invalid filename: a valid file extension is required",
			},
			toString: func(t testing.TB, envStruct any) string {
				// see previous
				return ""
			},
		},
		{
			name: "FailInvalidFilenameWrongFormat",
			args: args{
				configPath: "./testdata",
				filename:   "config",
				envStruct:  nil,
			},
			expected: expected{
				errString: "invalid filename: a filename of this format 'name.ext' is required",
			},
			toString: func(t testing.TB, envStruct any) string {
				// see previous
				return ""
			},
		},
		{
			name: "FailConfigDoesNotExist",
			args: args{
				configPath: "./testdata",
				filename:   "fakeconfig.json",
				envStruct:  nil,
			},
			expected: expected{
				errString: "fakeconfig.json: no such config file found in directory ./testdata",
			},
			toString: func(t testing.TB, envStruct any) string {
				// see previous
				return ""
			},
		},
		{
			// supposedtobejson.toml is a toml file, but I passed it the filename supposedtobejson.json instead, which doesn't exist
			// I don't know viper will recognize that it's the wrong extension, so I added some more descriptive errors for these cases
			name: "FailWrongExtensionType",
			args: args{
				configPath: "./testdata",
				filename:   "supposedtobejson.json",
				envStruct:  nil,
			},
			expected: expected{
				errString: "supposedtobejson.json: no such config file found in directory ./testdata",
			},
			toString: func(t testing.TB, envStruct any) string {
				// see previous
				return ""
			},
		},
		{
			name: "FailInvalidJSONData",
			args: args{
				configPath: "./testdata",
				filename:   "broken.json",
				envStruct:  nil,
			},
			expected: expected{
				errString: "While parsing config: unexpected end of JSON input",
			},
			toString: func(t testing.TB, envStruct any) string {
				// see previous
				return ""
			},
		},
		{
			name: "FailInvalidTOMLData",
			args: args{
				configPath: "./testdata",
				filename:   "broken.toml",
				envStruct:  nil,
			},
			expected: expected{
				errString: "While parsing config: toml: invalid character at start of key: {",
			},
			toString: func(t testing.TB, envStruct any) string {
				// see previous
				return ""
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := LoadFile(
				tt.args.configPath,
				tt.args.filename,
				&tt.args.envStruct,
			); err != nil {
				assert.EqualError(t, err, tt.expected.errString)
			} else {
				actual := tt.toString(t, &tt.args.envStruct)
				assert.JSONEq(t, tt.expected.envData, string(actual))
			}
		})
	}
}
