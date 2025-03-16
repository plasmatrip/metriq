package config

import (
	"crypto/rsa"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/caarlos0/env"
	"github.com/plasmatrip/metriq/internal/server/cert"
)

const (
	port               = "8080"
	host               = "localhost"
	storeinterval      = 300
	fileStoragePath    = "backup.dat"
	restore            = true
	retryInterval      = time.Second * 2
	startRetryInterval = time.Second * 1
	maxRetries         = 3
)

type Config struct {
	Host               string `env:"ADDRESS"`           // адрес сервера
	StoreInterval      int    `env:"STORE_INTERVAL"`    // интервал сохранения метрик
	FileStoragePath    string `env:"FILE_STORAGE_PATH"` // путь к файлу c метриками
	Restore            bool   `env:"RESTORE"`           // загружать ли сохраненные метрики
	DSN                string `env:"DATABASE_DSN"`      // подключение к бд
	Key                string `env:"KEY"`               // ключ для вычисления хэша по SHA256
	CryptoKeyPath      string `env:"CRYPTO_KEY"`        // путь к секретному ключу
	CryptoKey          *rsa.PrivateKey
	RetryInterval      time.Duration // увеличиваем интервал в сек между попытками повторного коннекта с бд
	StartRetryInterval time.Duration // начиниаем повторную попытку коннекта с бд через сек
	MaxRetries         int           // максимальное количество попыток повторного коннекта с бд
}

func NewConfig() (*Config, error) {
	cfg := &Config{
		RetryInterval:      retryInterval,
		StartRetryInterval: startRetryInterval,
		MaxRetries:         maxRetries,
	}

	cl := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

	// читаем переменные окружения, при ошибке прокидываем ее наверх
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("failed to read environment variable: %w", err)
	}

	var fHost string
	cl.StringVar(&fHost, "a", host+":"+port, "server address host:port")

	var fStoreInterval int
	cl.IntVar(&fStoreInterval, "i", storeinterval, "time interval in seconds for saving the metrics to a file")

	var fFileStoragePath string
	cl.StringVar(&fFileStoragePath, "f", fileStoragePath, "path to the file where metrics are saved")

	var fRestore bool
	cl.BoolVar(&fRestore, "r", restore, "whether to load saved metrics from a file or not")

	var fDSN string
	cl.StringVar(&fDSN, "d", "", "data source name to connect to the database")

	var fKey string
	cl.StringVar(&fKey, "k", "", "the key for calculating the hash using the SHA256 algorithm")

	var fCryptoKeyPath string
	cl.StringVar(&fCryptoKeyPath, "crypto-key", "", "the key for encrypting metrics")

	if err := cl.Parse(os.Args[1:]); err != nil {
		return nil, fmt.Errorf("failed to parse flags: %w", err)
	}

	if _, exist := os.LookupEnv("ADDRESS"); !exist {
		cfg.Host = fHost
	}

	if _, exist := os.LookupEnv("STORE_INTERVAL"); !exist {
		cfg.StoreInterval = fStoreInterval
	}

	if _, exist := os.LookupEnv("FILE_STORAGE_PATH"); !exist {
		cfg.FileStoragePath = fFileStoragePath
	}

	if _, exist := os.LookupEnv("RESTORE"); !exist {
		cfg.Restore = fRestore
	}

	if _, exist := os.LookupEnv("DATABASE_DSN"); !exist {
		cfg.DSN = fDSN
	}

	if _, exist := os.LookupEnv("KEY"); !exist {
		cfg.Key = fKey
	}

	if _, exist := os.LookupEnv("CRYPTO_KEY"); !exist {
		cfg.CryptoKeyPath = fCryptoKeyPath
	}

	if cfg.CryptoKey != nil {
		var err error
		cfg.CryptoKey, err = cert.LoadPrivateKey(cfg.CryptoKeyPath)
		if err != nil {
			return nil, fmt.Errorf("failed to get public key from cert: %w", err)
		}
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
