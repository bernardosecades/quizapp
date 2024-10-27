//go:build unit

package env

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetEnv(t *testing.T) {
	t.Parallel()

	type args struct {
		key      string
		fallback string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Environment variable is not set, should return fallback",
			args: args{
				key:      "BASE_URL",
				fallback: "http://lololo.com",
			},
			want: "http://lololo.com",
		},
		{
			name: "Environment variable is set, should return its value",
			args: args{
				key:      "BASE_URL",
				fallback: "http://fallback.com",
			},
			want: "http://lalala.com",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() { _ = os.Unsetenv(tt.args.key) }()
			if tt.want != tt.args.fallback {
				_ = os.Setenv(tt.args.key, tt.want)
			}
			assert.Equal(t, GetEnv(tt.args.key, tt.args.fallback), tt.want, tt.name)
		})
	}
}
