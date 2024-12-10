package config

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/caarlos0/env/v6"
)

const (
	port               = "8080"
	host               = "localhost"
	pollInterval       = 2
	reportInterval     = 10
	clientTimeout      = time.Second * 5
	retryInterval      = time.Second * 2
	startRetryInterval = time.Second * 1
	maxRetries         = 3
)

type Config struct {
	Host               string        `env:"ADDRESS"`         // адрес сервера
	PollInterval       int           `env:"POLL_INTERVAL"`   // интервал в сек обновления метрик
	ReportInterval     int           `env:"REPORT_INTERVAL"` // интервал в сек отправки метрик на сервер
	Key                string        `env:"KEY"`             // ключ для вычисления хэша по SHA256
	ClientTimeout      time.Duration // таймаут для http клиента
	RetryInterval      time.Duration // увеличиваем интервал в сек между попытками повторной отправки метрик на сервер
	StartRetryInterval time.Duration // начиниаем повторную отправку через сек
	MaxRetries         int           // максимальное количество попыток повторной отправки метрик на сервер
}

func NewConfig() (*Config, error) {
	cfg := &Config{
		ClientTimeout:      clientTimeout,
		RetryInterval:      retryInterval,
		StartRetryInterval: startRetryInterval,
		MaxRetries:         maxRetries,
	}

	cl := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

	// читаем переменную окружения, при ошибке прокидываем ее наверх
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("failed to read environment variable: %w", err)
	}

	// проверяем флаги

	var fHost string
	cl.StringVar(&fHost, "a", host+":"+port, "server address host:port")

	var fPollInterval int
	cl.IntVar(&fPollInterval, "p", pollInterval, "metrics reporting interval")

	var fReportInterval int
	cl.IntVar(&fReportInterval, "r", reportInterval, "metrics polling frequency")

	var fKey string
	cl.StringVar(&fKey, "k", "", "the key for calculating the hash using the SHA256 algorithm")

	// при ошибке парсинга прокидываем ошибку наверх
	if err := cl.Parse(os.Args[1:]); err != nil {
		return nil, fmt.Errorf("failed to parse flags: %w", err)
	}

	// если переменная есть парсим адрес, если порт задан не числом прокидываем ошибку наверх
	// if len(cfg.Host) != 0 {
	// 	if err := parseAddress(cfg); err != nil {
	// 		return nil, fmt.Errorf("port parsing error: %w", err)
	// 	}

	// 	if cfg.PollInterval <= 0 {
	// 		cfg.PollInterval = pollInterval
	// 	}

	// 	if cfg.ReportInterval <= 0 {
	// 		cfg.ReportInterval = reportInterval
	// 	}

	// 	return cfg, nil
	// }

	if _, exist := os.LookupEnv("ADDRESS"); !exist {
		cfg.Host = fHost
	}

	if _, exist := os.LookupEnv("POLL_INTERVAL"); !exist {
		cfg.PollInterval = fPollInterval
	}

	if _, exist := os.LookupEnv("REPORT_INTERVAL"); !exist {
		cfg.ReportInterval = fReportInterval
	}

	if _, exist := os.LookupEnv("KEY"); !exist {
		cfg.Key = fKey
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
