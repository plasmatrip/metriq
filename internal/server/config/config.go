package config

import (
	"crypto/rsa"
	"encoding/json"
	"flag"
	"fmt"
	"net"
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
	ConfFile           string          `env:"CONFIG"`                               // путь к конфигурационному File
	Host               string          `env:"ADDRESS" json:"address"`               // адрес сервера
	StoreInterval      int             `env:"STORE_INTERVAL" json:"store_interval"` // интервал сохранения метрик
	FileStoragePath    string          `env:"FILE_STORAGE_PATH" json:"store_file"`  // путь к файлу c метриками
	Restore            bool            `env:"RESTORE" json:"restore"`               // загружать ли сохраненные метрики
	DSN                string          `env:"DATABASE_DSN" json:"database_dsn"`     // подключение к бд
	Key                string          `env:"KEY"`                                  // ключ для вычисления хэша по SHA256
	CryptoKeyPath      string          `env:"CRYPTO_KEY" json:"crypto_key"`         // путь к секретному ключу
	TrustedSubnet      string          `env:"TRUSTED_SUBNET" json:"trusted_subnet"` // сеть с которой разрешен доступ
	TrustedSubnetCIDR  *net.IPNet      // сеть с которой разрешен доступ
	CryptoKey          *rsa.PrivateKey // секретный ключ
	RetryInterval      time.Duration   // увеличиваем интервал в сек между попытками повторного коннекта с бд
	StartRetryInterval time.Duration   // начиниаем повторную попытку коннекта с бд через сек
	MaxRetries         int             // максимальное количество попыток повторного коннекта с бд
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

	var fConfig string
	cl.StringVar(&fConfig, "c", "", "path to the configuration file")
	cl.StringVar(&fConfig, "config", "", "path to the configuration file")

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

	var fTrustedSubnet string
	cl.StringVar(&fTrustedSubnet, "t", "", "the subnet from which access is allowed")

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

	if _, exist := os.LookupEnv("ADDRESS"); !exist && fHost != "" {
		cfg.Host = fHost
	}

	if _, exist := os.LookupEnv("STORE_INTERVAL"); !exist && fStoreInterval != 0 {
		cfg.StoreInterval = fStoreInterval
	}

	if _, exist := os.LookupEnv("FILE_STORAGE_PATH"); !exist && fFileStoragePath != "" {
		cfg.FileStoragePath = fFileStoragePath
	}

	if _, exist := os.LookupEnv("RESTORE"); !exist && fRestore {
		cfg.Restore = fRestore
	}

	if _, exist := os.LookupEnv("DATABASE_DSN"); !exist && fDSN != "" {
		cfg.DSN = fDSN
	}

	if _, exist := os.LookupEnv("KEY"); !exist {
		cfg.Key = fKey
	}

	if _, exist := os.LookupEnv("CRYPTO_KEY"); !exist && fCryptoKeyPath != "" {
		cfg.CryptoKeyPath = fCryptoKeyPath
	}

	if _, exist := os.LookupEnv("TRUSTED_SUBNET"); !exist && fTrustedSubnet != "" {
		cfg.TrustedSubnet = fTrustedSubnet
	}

	if cfg.TrustedSubnet != "" {
		var err error
		_, cfg.TrustedSubnetCIDR, err = net.ParseCIDR(cfg.TrustedSubnet)
		if err != nil {
			return nil, fmt.Errorf("failed to parse trusted subnet: %w", err)
		}
	}

	if cfg.CryptoKeyPath != "" {
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
