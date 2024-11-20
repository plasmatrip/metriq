package config

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/caarlos0/env/v6"
)

const (
	port           = "8080"
	host           = "localhost"
	pollInterval   = 2
	reportInterval = 10
)

type Config struct {
	Host           string `env:"ADDRESS"`
	PollInterval   int    `env:"POLL_INTERVAL"`
	ReportInterval int    `env:"REPORT_INTERVAL"`
}

func NewConfig() (*Config, error) {
	cfg := new(Config)

	cl := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

	// читаем переменную окружения, при ошибке прокидываем ее наверх
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("failed to read environment variable: %w", err)
	}

	// если переменная есть парсим адрес, если порт задан не числом прокидываем ошибку наверх
	if _, exist := os.LookupEnv("ADDRESS"); !exist {
		cl.StringVar(&cfg.Host, "a", host+":"+port, "server address host:port")
	}

	if _, exist := os.LookupEnv("POLL_INTERVAL"); !exist {
		cl.IntVar(&cfg.PollInterval, "p", pollInterval, "metrics reporting interval")
	}

	if _, exist := os.LookupEnv("REPORT_INTERVAL"); !exist {
		cl.IntVar(&cfg.ReportInterval, "r", reportInterval, "metrics polling frequency")
	}

	// при ошибке парсинга прокидываем ошибку наверх
	if err := cl.Parse(os.Args[1:]); err != nil {
		return nil, fmt.Errorf("failed to parse flags: %w", err)
	}

	if err := parseAddress(cfg); err != nil {
		return nil, fmt.Errorf("port parsing error: %w", err)
	}

	if cfg.PollInterval <= 0 {
		cfg.PollInterval = pollInterval
	}

	if cfg.ReportInterval <= 0 {
		cfg.ReportInterval = reportInterval
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
