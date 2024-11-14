package server

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/caarlos0/env/v6"
)

const (
	port = "8080"
	host = "localhost"

	UpdateURLLen = 5
	ValueURLLen  = 4

	RequestTypePos  = 2
	RequestNamePos  = 3
	RequestValuePos = 4

	Gauge   = "gauge"
	Counter = "counter"

	PollCount = "PollCount"
)

type Config struct {
	Host string `env:"ADDRESS"`
}

func NewConfig() (*Config, error) {
	config := new(Config)

	// читаем переменную окружения, при ошибке прокидываем ее наверх
	if err := env.Parse(config); err != nil {
		return nil, fmt.Errorf("failed to read environment variable: %w", err)
	}

	// если переменная есть парсим адрес, если порт задан не числом прокидываем ошибку наверх
	if len(config.Host) != 0 {
		if err := parseAddress(config); err != nil {
			return nil, fmt.Errorf("port parsing error: %w", err)
		}

		return config, nil
	}

	//проверяем флаги
	cl := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	cl.StringVar(&config.Host, "a", host+":"+port, "Server address host:port")

	// ошибке парсинга прокидываем ошибку наверх
	if err := cl.Parse(os.Args[1:]); err != nil {
		return nil, fmt.Errorf("failed to parse flags: %w", err)
	}

	if err := parseAddress(config); err != nil {
		return nil, fmt.Errorf("port parsing error: %w", err)
	}

	return config, nil
}

func parseAddress(config *Config) error {
	args := strings.Split(config.Host, ":")
	if len(args) == 2 {
		if len(args[0]) == 0 || len(args[1]) == 0 {
			config.Host = host + ":" + port
			return nil
		}

		_, err := strconv.ParseInt(args[1], 10, 64)
		return err
	}
	config.Host = host + ":" + port
	return nil
}
