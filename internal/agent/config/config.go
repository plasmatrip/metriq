// This package contains the configuration for the metrics collection agent.
// The agent is responsible for collecting various metrics about the application
// and sending them to the server. The configuration is read from environment
// variables, flags and a configuration file. The configuration provides the
// server address, the polling interval for collecting metrics and the path to
// the file where metrics are saved. The agent uses a context to listen for
// termination signals and sets up goroutines for collecting application
// metrics, collecting system metrics using gopsutil and sending metrics to the
// server. The agent gracefully shuts down all goroutines and exits when it
// receives a termination signal.
package config

import (
	"crypto/rsa"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/plasmatrip/metriq/internal/agent/cert"
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
	rateLimit          = 5
)

type Config struct {
	ConfFile           string `env:"CONFIG"`          // путь к конфигурационному File
	Host               string `env:"ADDRESS"`         // адрес сервера
	PollInterval       int    `env:"POLL_INTERVAL"`   // интервал в сек обновления метрик
	ReportInterval     int    `env:"REPORT_INTERVAL"` // интервал в сек отправки метрик на сервер
	Key                string `env:"KEY"`             // ключ для вычисления хэша по SHA256
	RateLimit          int    `env:"RATE_LIMIT"`      // количество одновременно исходящих запросов на сервер
	CryptoKeyPath      string `env:"CRYPTO_KEY"`      // ауть к сертификату
	CryptoKey          *rsa.PublicKey
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
	var fConfig string
	cl.StringVar(&fConfig, "c", "", "path to the configuration file")
	cl.StringVar(&fConfig, "config", "", "path to the configuration file")

	var fHost string
	cl.StringVar(&fHost, "a", host+":"+port, "server address host:port")

	var fPollInterval int
	cl.IntVar(&fPollInterval, "p", pollInterval, "metrics reporting interval")

	var fReportInterval int
	cl.IntVar(&fReportInterval, "r", reportInterval, "metrics polling frequency")

	var fRateLimit int
	cl.IntVar(&fRateLimit, "l", rateLimit, "maximum number of workers")

	var fKey string
	cl.StringVar(&fKey, "k", "", "the key for calculating the hash using the SHA256 algorithm")

	var fCryptoKeyPath string
	cl.StringVar(&fCryptoKeyPath, "crypto-key", "", "the key for encrypting metrics")

	// при ошибке парсинга прокидываем ошибку наверх
	if err := cl.Parse(os.Args[1:]); err != nil {
		return nil, fmt.Errorf("failed to parse flags: %w", err)
	}

	if _, exist := os.LookupEnv("CONFIG"); !exist {
		cfg.ConfFile = fConfig
	}

	// читаем конфигурационный файл
	if cfg.ConfFile != "" {
		data, err := os.ReadFile(cfg.ConfFile)
		if err != nil {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}

		err = json.Unmarshal(data, &cfg)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal config file: %w", err)
		}
	}

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

	if _, exist := os.LookupEnv("RATE_LIMIT"); !exist {
		cfg.RateLimit = fRateLimit
	}

	if _, exist := os.LookupEnv("CRYPTO_KEY"); !exist {
		cfg.CryptoKeyPath = fCryptoKeyPath
	}

	if cfg.CryptoKey != nil {
		var err error
		cfg.CryptoKey, err = cert.GetPublicKeyFromCert(cfg.CryptoKeyPath)
		if err != nil {
			return nil, fmt.Errorf("failed to get public key from cert: %w", err)
		}
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
