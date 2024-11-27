package config

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/caarlos0/env"
)

const (
	port            = "8080"
	host            = "localhost"
	storeinterval   = 300
	fileStoragePath = "backup.dat"
	restore         = true
)

type Config struct {
	Host            string `env:"ADDRESS"`
	StoreInterval   int    `env:"STORE_INTERVAL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	Restore         bool   `env:"RESTORE"`
	DSN             string `env:"DATABASE_DSN"`
}

func NewConfig() (*Config, error) {
	cfg := new(Config)

	cl := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

	// читаем переменные окружения, при ошибке прокидываем ее наверх
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("failed to read environment variable: %w", err)
	}

	// если переменная есть парсим адрес, если порт задан не числом прокидываем ошибку наверх
	if _, exist := os.LookupEnv("ADDRESS"); !exist {
		cl.StringVar(&cfg.Host, "a", host+":"+port, "Server address host:port")
	}

	if _, exist := os.LookupEnv("STORE_INTERVAL"); !exist {
		cl.IntVar(&cfg.StoreInterval, "i", storeinterval, "Time interval in seconds for saving the metrics to a file")
	}

	if _, exist := os.LookupEnv("FILE_STORAGE_PATH"); !exist {
		cl.StringVar(&cfg.FileStoragePath, "f", fileStoragePath, "Path to the file where metrics are saved")
	}

	if _, exist := os.LookupEnv("RESTORE"); !exist {
		cl.BoolVar(&cfg.Restore, "r", restore, "Whether to load saved metrics from a file or not")
	}

	if _, exist := os.LookupEnv("DATABASE_DSN"); !exist {
		cl.StringVar(&cfg.DSN, "d", "", "Data source name to connect to the database")
	}

	if err := cl.Parse(os.Args[1:]); err != nil {
		return nil, fmt.Errorf("failed to parse flags: %w", err)
	}

	if err := parseAddress(cfg); err != nil {
		return nil, fmt.Errorf("port parsing error: %w", err)
	}

	return cfg, nil
}

func parseAddress(cfg *Config) error {
	args := strings.Split(cfg.Host, ":")
	if len(args) == 2 {
		if len(args[0]) == 0 || len(args[1]) == 0 {
			cfg.Host = host + ":" + port
			return nil
		}

		_, err := strconv.ParseInt(args[1], 10, 64)
		return err
	}
	cfg.Host = host + ":" + port
	return nil
}

// package config

// import (
// 	"flag"
// 	"fmt"
// 	"os"
// 	"strconv"
// 	"strings"

// 	"github.com/caarlos0/env"
// )

// const (
// 	port            = "8080"
// 	host            = "localhost"
// 	storeinterval   = 300
// 	fileStoragePath = "backup.dat"
// 	restore         = true
// )

// type Config struct {
// 	Host            string `env:"ADDRESS"`
// 	StoreInterval   int    `env:"STORE_INTERVAL"`
// 	FileStoragePath string `env:"FILE_STORAGE_PATH"`
// 	Restore         bool   `env:"RESTORE"`
// }

// func NewConfig() (*Config, error) {
// 	cfg := new(Config)

// 	cl := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

// 	// читаем переменные окружения, при ошибке прокидываем ее наверх
// 	if err := env.Parse(cfg); err != nil {
// 		return nil, fmt.Errorf("failed to read environment variable: %w", err)
// 	}

// 	// если переменная есть парсим адрес, если порт задан не числом прокидываем ошибку наверх
// 	if _, exist := os.LookupEnv("ADDRESS"); !exist {
// 		cl.StringVar(&cfg.Host, "a", host+":"+port, "Server address host:port")
// 	}
// 	if err := parseAddress(cfg); err != nil {
// 		return nil, fmt.Errorf("port parsing error: %w", err)
// 	}

// 	if _, exist := os.LookupEnv("STORE_INTERVAL"); !exist {
// 		cl.IntVar(&cfg.StoreInterval, "i", storeinterval, "Time interval in seconds for saving the metrics to a file")
// 	}

// 	if _, exist := os.LookupEnv("FILE_STORAGE_PATH"); !exist {
// 		cl.StringVar(&cfg.FileStoragePath, "f", fileStoragePath, "Path to the file where metrics are saved")
// 	}

// 	if _, exist := os.LookupEnv("RESTORE"); !exist {
// 		cl.BoolVar(&cfg.Restore, "r", restore, "Whether to load saved metrics from a file or not")
// 	}

// 	if err := cl.Parse(os.Args[1:]); err != nil {
// 		return nil, fmt.Errorf("failed to parse flags: %w", err)
// 	}

// 	return cfg, nil
// }

// func parseAddress(config *Config) error {
// 	args := strings.Split(config.Host, ":")
// 	if len(args) == 2 {
// 		if len(args[0]) == 0 || len(args[1]) == 0 {
// 			config.Host = host + ":" + port
// 			return nil
// 		}

// 		_, err := strconv.ParseInt(args[1], 10, 64)
// 		return err
// 	}
// 	config.Host = host + ":" + port
// 	return nil
// }
