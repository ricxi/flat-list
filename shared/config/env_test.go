package config

import (
	"testing"
)

func TestEnvConfigErr_Error(t *testing.T) {
	type fields struct {
		missingEnvs []string
		invalidEnvs []string
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
				invalidEnvs: []string{},
			},
			want: "",
		},
		{
			name: "NoNilDereference",
			fields: fields{
				missingEnvs: nil,
				invalidEnvs: nil,
			},
			want: "",
		},
		{
			name: "OnlyOneMissingEnvError",
			fields: fields{
				missingEnvs: []string{"PORT"},
				invalidEnvs: nil,
			},
			want: "missing [PORT] environment variable(s)",
		},
		{
			name: "OnlyMissingEnvErrors",
			fields: fields{
				missingEnvs: []string{"PORT", "HOST", "DBNAME"},
				invalidEnvs: nil,
			},
			want: "missing [PORT HOST DBNAME] environment variable(s)",
		},
		{
			name: "OnlyOneInvalidEnvError",
			fields: fields{
				missingEnvs: nil,
				invalidEnvs: []string{"PORT"},
			},
			want: "invalid [PORT] environment variable(s)",
		},
		{
			name: "OnlyInvalidEnvErrors",
			fields: fields{
				missingEnvs: nil,
				invalidEnvs: []string{"PORT HOST DB_URI"},
			},
			want: "invalid [PORT HOST DB_URI] environment variable(s)",
		},
		{
			name: "BothInvalidAndMissingEnvErrors",
			fields: fields{
				missingEnvs: []string{"PORT", "HOST"},
				invalidEnvs: []string{"DB_NAME", "SMTP_ADDR", "DB_URI"},
			},
			want: "missing [PORT HOST] and invalid [DB_NAME SMTP_ADDR DB_URI] environment variable(s)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &EnvConfigErr{
				missingEnvs: tt.fields.missingEnvs,
				invalidEnvs: tt.fields.invalidEnvs,
			}
			if got := m.Error(); got != tt.want {
				t.Errorf("got EnvConfigErr.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}
