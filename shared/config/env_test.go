package config_test

import (
	"testing"

	"github.com/ricxi/flat-list/shared/config"
)

func TestLoadEnvs(t *testing.T) {
	type args struct {
		envKeys   []string
		envValues []string
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]string
		wantErr bool
	}{
		{
			name: "SuccessLoadTwoEnvs",
			args: args{
				envKeys:   []string{"PORT", "HOST"},
				envValues: []string{"5000", "127.0.0.1"},
			},
			want:    map[string]string{"PORT": "5000", "HOST": "127.0.0.1"},
			wantErr: false,
		},
		{
			name: "FailMissingOneEnv",
			args: args{
				envKeys:   []string{"PORT", "HOST"},
				envValues: []string{"5000"},
			},
			want:    map[string]string{},
			wantErr: true,
		},
		{
			name: "FailNoEnvs",
			args: args{
				envKeys:   nil,
				envValues: nil,
			},
			want:    map[string]string{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for k, v := range tt.args.envValues {
				t.Setenv(tt.args.envKeys[k], v)
			}

			got, err := config.LoadEnvs(tt.args.envKeys...)
			if (err != nil) != tt.wantErr {
				t.Errorf("got LoadEnvs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !compareEnvMaps(got, tt.want) {
				t.Errorf("got LoadEnvs() error = %v, wantErr %v", err, tt.wantErr)
				t.Errorf("got LoadEnvs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func compareEnvMaps(m1, m2 map[string]string) bool {
	if len(m1) != len(m2) {
		return false
	}

	for k, v := range m1 {
		if m2[k] != v {
			return false
		}
	}

	return true
}

func Test_envMap_ValidateAsInt(t *testing.T) {
	type args struct {
		envkey   string
		envValue string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "SuccessfulIntConversion",
			args: args{
				envkey:   "PORT",
				envValue: "5000",
			},
			wantErr: false,
		},
		{
			name: "FailIntConversion",
			args: args{
				envkey:   "PORT",
				envValue: "STRINGSHOULDFAIL5000",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv(tt.args.envkey, tt.args.envValue)

			em, err := config.LoadEnvs(tt.args.envkey)
			if err != nil {
				t.Fatal("got error, but did not expect one for this test")
				return
			}

			if err := em.ValidateAsInt(tt.args.envkey); (err != nil) != tt.wantErr {
				t.Errorf("envMap.ValidateAsInt() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
