package agent

import (
	"flag"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfig_ParseAddress(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  string
	}{
		{
			name:  "Valid address",
			value: "server.com:8989",
			want:  "server.com:8989",
		},
		{
			name:  "Empty port number",
			value: "server.com:",
			want:  "localhost:8080",
		},
		{
			name:  "Empty host name",
			value: ":8574",
			want:  "localhost:8080",
		},
		{
			name:  "Empty address",
			value: "",
			want:  "localhost:8080",
		},
		{
			name:  "Only colon",
			value: ":",
			want:  "localhost:8080",
		},
	}
	config := new(Config)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			config.Host = test.value
			parseAddress(config)
			assert.Equal(t, test.want, config.Host)
		})
	}
}

func TestConfig_NewConfig(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want Config
	}{
		{
			name: "Valid config",
			args: []string{"-a", "server.com:8585"},
			want: Config{Host: "server.com:8585", PollInterval: 2, ReportInterval: 10},
		},
	}
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fs.String("a", "localhost:8080", "Server address host:port")
			fs.Parse(test.args)
			config, err := NewConfig()
			require.NoError(t, err)
			assert.Equal(t, test.want, *config)
		})
	}
}
