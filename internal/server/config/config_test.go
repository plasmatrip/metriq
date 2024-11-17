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
			env:     map[string]string{"ADDRESS": "server.com:8585"},
			want:    Config{Host: "server.com:8585"},
			errWant: false,
		},
		{
			name:    "Invalid port type",
			env:     map[string]string{"ADDRESS": "server.com:ttt"},
			want:    Config{},
			errWant: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
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

// func TestConfig_NewConfig_Flags(t *testing.T) {
// 	tests := []struct {
// 		name    string
// 		args    []string
// 		want    Config
// 		errWant bool
// 	}{
// 		{
// 			name:    "Valid config",
// 			args:    []string{},
// 			want:    Config{Host: "localhost:8080", StoreInterval: 300, FileStoragePath: "./data/backup.dat", Restore: true},
// 			errWant: false,
// 		},
// 		{
// 			name:    "Valid config",
// 			args:    []string{"-a", "server.com:8585"},
// 			want:    Config{Host: "server.com:8585", StoreInterval: 300, FileStoragePath: "./data/backup.dat", Restore: true},
// 			errWant: false,
// 		},
// 		{
// 			name:    "Empty port",
// 			args:    []string{"-a", "server.com:"},
// 			want:    Config{Host: "localhost:8080", StoreInterval: 300, FileStoragePath: "./data/backup.dat", Restore: true},
// 			errWant: false,
// 		},
// 		{
// 			name:    "Empty host name",
// 			args:    []string{"-a", ":8585"},
// 			want:    Config{Host: "localhost:8080", StoreInterval: 300, FileStoragePath: "./data/backup.dat", Restore: true},
// 			errWant: false,
// 		},
// 		{
// 			name:    "Empty address",
// 			args:    []string{"-a", ""},
// 			want:    Config{Host: "localhost:8080", StoreInterval: 300, FileStoragePath: "./data/backup.dat", Restore: true},
// 			errWant: false,
// 		},
// 		{
// 			name:    "Only colon",
// 			args:    []string{"-a", ":"},
// 			want:    Config{Host: "localhost:8080", StoreInterval: 300, FileStoragePath: "./data/backup.dat", Restore: true},
// 			errWant: false,
// 		},
// 		{
// 			name:    "Extra flag",
// 			args:    []string{"-a", "localhost:8080", "-s", "0"},
// 			want:    Config{},
// 			errWant: true,
// 		},
// 		{
// 			name:    "Invalid port type",
// 			args:    []string{"-a", "server.com:ddd"},
// 			want:    Config{},
// 			errWant: true,
// 		},
// 	}

// 	for _, test := range tests {
// 		t.Run(test.name, func(t *testing.T) {
// 			os.Args = []string{os.Args[0]}
// 			if len(test.args) > 0 {
// 				os.Args = append(os.Args, test.args...)
// 			}
// 			config, err := NewConfig()
// 			if test.errWant {
// 				require.Error(t, err)
// 				return
// 			}
// 			require.NoError(t, err)
// 			assert.Equal(t, test.want, *config)
// 		})
// 	}
// }
