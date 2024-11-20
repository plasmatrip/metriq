package config

import (
	"os"
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

func TestConfig_NewConfig_env(t *testing.T) {
	tests := []struct {
		name    string
		env     map[string]string
		want    Config
		errWant bool
	}{
		{
			name:    "Valid config",
			env:     map[string]string{"ADDRESS": "server.com:8585", "POLL_INTERVAL": "2", "REPORT_INTERVAL": "10"},
			want:    Config{Host: "server.com:8585", PollInterval: 2, ReportInterval: 10},
			errWant: false,
		},
		{
			name:    "Valid config",
			env:     map[string]string{"ADDRESS": "server.com:8585", "POLL_INTERVAL": "5", "REPORT_INTERVAL": "10"},
			want:    Config{Host: "server.com:8585", PollInterval: 5, ReportInterval: 10},
			errWant: false,
		},
		{
			name:    "Valid config",
			env:     map[string]string{"ADDRESS": "server.com:8585", "POLL_INTERVAL": "2", "REPORT_INTERVAL": "15"},
			want:    Config{Host: "server.com:8585", PollInterval: 2, ReportInterval: 15},
			errWant: false,
		},
		{
			name:    "Invalid port type",
			env:     map[string]string{"ADDRESS": "server.com:ttt", "POLL_INTERVAL": "2", "REPORT_INTERVAL": "10"},
			want:    Config{},
			errWant: true,
		},
		{
			name:    "Invalid poll interval type",
			env:     map[string]string{"ADDRESS": "server.com:8585", "POLL_INTERVAL": "p", "REPORT_INTERVAL": "10"},
			want:    Config{},
			errWant: true,
		},
		{
			name:    "Invalid report interval type",
			env:     map[string]string{"ADDRESS": "server.com:8585", "POLL_INTERVAL": "2", "REPORT_INTERVAL": "r"},
			want:    Config{},
			errWant: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			os.Args = []string{os.Args[0]}
			for k, v := range test.env {
				os.Setenv(k, v)
			}
			config, err := NewConfig()
			if test.errWant {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.want, *config)
		})
	}

	os.Clearenv()
}

func TestConfig_NewConfig_Flags(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		want    Config
		errWant bool
	}{
		{
			name:    "Valid config",
			args:    []string{},
			want:    Config{Host: "localhost:8080", PollInterval: 2, ReportInterval: 10},
			errWant: false,
		},
		{
			name:    "Valid config",
			args:    []string{"-a", "server.com:8585"},
			want:    Config{Host: "server.com:8585", PollInterval: 2, ReportInterval: 10},
			errWant: false,
		},
		{
			name:    "Valid config",
			args:    []string{"-p", "2"},
			want:    Config{Host: "localhost:8080", PollInterval: 2, ReportInterval: 10},
			errWant: false,
		},
		{
			name:    "Valid config",
			args:    []string{"-r", "10"},
			want:    Config{Host: "localhost:8080", PollInterval: 2, ReportInterval: 10},
			errWant: false,
		},
		{
			name:    "Valid config",
			args:    []string{"-a", "server.com:8585", "-p", "2"},
			want:    Config{Host: "server.com:8585", PollInterval: 2, ReportInterval: 10},
			errWant: false,
		},
		{
			name:    "Valid config",
			args:    []string{"-a", "server.com:8585", "-p", "2", "-r", "10"},
			want:    Config{Host: "server.com:8585", PollInterval: 2, ReportInterval: 10},
			errWant: false,
		},
		{
			name:    "Empty port",
			args:    []string{"-a", "server.com:", "-p", "2", "-r", "10"},
			want:    Config{Host: "localhost:8080", PollInterval: 2, ReportInterval: 10},
			errWant: false,
		},
		{
			name:    "Empty host name",
			args:    []string{"-a", ":8585", "-p", "2", "-r", "10"},
			want:    Config{Host: "localhost:8080", PollInterval: 2, ReportInterval: 10},
			errWant: false,
		},
		{
			name:    "Empty address",
			args:    []string{"-a", "", "-p", "2", "-r", "10"},
			want:    Config{Host: "localhost:8080", PollInterval: 2, ReportInterval: 10},
			errWant: false,
		},
		{
			name:    "Only colon",
			args:    []string{"-a", ":", "-p", "2", "-r", "10"},
			want:    Config{Host: "localhost:8080", PollInterval: 2, ReportInterval: 10},
			errWant: false,
		},
		{
			name:    "Extra flag",
			args:    []string{"-a", "localhost:8080", "-p", "2", "-r", "10", "-s", "0"},
			want:    Config{},
			errWant: true,
		},
		{
			name:    "Incorrect value",
			args:    []string{"-r", "r"},
			want:    Config{},
			errWant: true,
		},
		{
			name:    "Negative poll interval",
			args:    []string{"-p", "-2"},
			want:    Config{Host: "localhost:8080", PollInterval: 2, ReportInterval: 10},
			errWant: false,
		},
		{
			name:    "Negative report interval",
			args:    []string{"-r", "-10"},
			want:    Config{Host: "localhost:8080", PollInterval: 2, ReportInterval: 10},
			errWant: false,
		},
		{
			name:    "Invalid port type",
			args:    []string{"-a", "server.com:ddd", "-p", "2", "-r", "10"},
			want:    Config{},
			errWant: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			os.Args = []string{os.Args[0]}
			if len(test.args) > 0 {
				os.Args = append(os.Args, test.args...)
			}
			config, err := NewConfig()
			if test.errWant {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.want, *config)
		})
	}
}
