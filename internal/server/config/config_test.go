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
			name:    "Valid server address",
			env:     map[string]string{"ADDRESS": "server.com:8585"},
			want:    Config{Host: "server.com:8585", StoreInterval: 300, FileStoragePath: "backup.dat", Restore: true, RetryInterval: 2000000000, StartRetryInterval: 1000000000, MaxRetries: 3},
			errWant: false,
		},
		{
			name:    "Invalid port type",
			env:     map[string]string{"ADDRESS": "server.com:ttt"},
			want:    Config{},
			errWant: true,
		},
		{
			name:    "Valid store interval",
			env:     map[string]string{"STORE_INTERVAL": "100"},
			want:    Config{Host: "localhost:8080", StoreInterval: 100, FileStoragePath: "backup.dat", Restore: true, RetryInterval: 2000000000, StartRetryInterval: 1000000000, MaxRetries: 3},
			errWant: false,
		},
		{
			name:    "Invalid store interval",
			env:     map[string]string{"STORE_INTERVAL": "ttt"},
			want:    Config{},
			errWant: true,
		},
		{
			name:    "Valid file storage",
			env:     map[string]string{"FILE_STORAGE_PATH": "file.dat"},
			want:    Config{Host: "localhost:8080", StoreInterval: 300, FileStoragePath: "file.dat", Restore: true, RetryInterval: 2000000000, StartRetryInterval: 1000000000, MaxRetries: 3},
			errWant: false,
		},
		{
			name:    "Valid restore",
			env:     map[string]string{"RESTORE": "false"},
			want:    Config{Host: "localhost:8080", StoreInterval: 300, FileStoragePath: "backup.dat", Restore: false, RetryInterval: 2000000000, StartRetryInterval: 1000000000, MaxRetries: 3},
			errWant: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			os.Clearenv()
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
			want:    Config{Host: "localhost:8080", StoreInterval: 300, FileStoragePath: "backup.dat", Restore: true, RetryInterval: 2000000000, StartRetryInterval: 1000000000, MaxRetries: 3},
			errWant: false,
		},
		{
			name:    "Valid config",
			args:    []string{"-a", "server.com:8585"},
			want:    Config{Host: "server.com:8585", StoreInterval: 300, FileStoragePath: "backup.dat", Restore: true, RetryInterval: 2000000000, StartRetryInterval: 1000000000, MaxRetries: 3},
			errWant: false,
		},
		{
			name:    "Empty port",
			args:    []string{"-a", "server.com:"},
			want:    Config{Host: "localhost:8080", StoreInterval: 300, FileStoragePath: "backup.dat", Restore: true, RetryInterval: 2000000000, StartRetryInterval: 1000000000, MaxRetries: 3},
			errWant: false,
		},
		{
			name:    "Empty host name",
			args:    []string{"-a", ":8585"},
			want:    Config{Host: "localhost:8080", StoreInterval: 300, FileStoragePath: "backup.dat", Restore: true, RetryInterval: 2000000000, StartRetryInterval: 1000000000, MaxRetries: 3},
			errWant: false,
		},
		{
			name:    "Empty address",
			args:    []string{"-a", ""},
			want:    Config{Host: "localhost:8080", StoreInterval: 300, FileStoragePath: "backup.dat", Restore: true, RetryInterval: 2000000000, StartRetryInterval: 1000000000, MaxRetries: 3},
			errWant: false,
		},
		{
			name:    "Only colon",
			args:    []string{"-a", ":"},
			want:    Config{Host: "localhost:8080", StoreInterval: 300, FileStoragePath: "backup.dat", Restore: true, RetryInterval: 2000000000, StartRetryInterval: 1000000000, MaxRetries: 3},
			errWant: false,
		},
		{
			name:    "Extra flag",
			args:    []string{"-a", "localhost:8080", "-s", "0"},
			want:    Config{},
			errWant: true,
		},
		{
			name:    "Invalid port type",
			args:    []string{"-a", "server.com:ddd"},
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
