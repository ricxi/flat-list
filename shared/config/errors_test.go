package config

import (
	"testing"
)

func TestEnvConfigErr_Error(t *testing.T) {
	type fields struct {
		missingEnvs []string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "NoErrors",
			fields: fields{
				missingEnvs: []string{},
			},
			want: "",
		},
		{
			name: "NoNilDereference",
			fields: fields{
				missingEnvs: nil,
			},
			want: "",
		},
		{
			name: "OnlyOneMissingEnvError",
			fields: fields{
				missingEnvs: []string{"PORT"},
			},
			want: "missing [PORT] environment variable(s)",
		},
		{
			name: "MissingEnvErrors",
			fields: fields{
				missingEnvs: []string{"PORT", "HOST", "DBNAME"},
			},
			want: "missing [PORT HOST DBNAME] environment variable(s)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &EnvConfigErr{
				missingEnvs: tt.fields.missingEnvs,
			}
			if got := m.Error(); got != tt.want {
				t.Errorf("got EnvConfigErr.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}
