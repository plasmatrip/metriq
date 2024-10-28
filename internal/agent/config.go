package agent

import (
	"flag"
	"fmt"
	"os"
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

func NewConfig() *Config {
	config := new(Config)

	//читаем переменную окружения, при ошибке выходим из программы
	err := env.Parse(config)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	//если переменная есть парсим адрес
	if len(config.Host) != 0 {
		parseAddress(config)

		return config
	}

	//проверяем флаги
	flag.StringVar(&config.Host, "a", "localhost:8080", "Server address host:port")
	flag.IntVar(&config.PollInterval, "p", pollInterval, "metrics reporting interval")
	flag.IntVar(&config.ReportInterval, "r", reportInterval, "metrics polling frequency")
	flag.Parse()

	if len(flag.Args()) > 0 {
		fmt.Println("Unknown flag(s): ", flag.Args())
		os.Exit(1)
	}

	parseAddress(config)

	if config.PollInterval <= 0 {
		config.PollInterval = pollInterval
	}

	if config.ReportInterval <= 0 {
		config.ReportInterval = reportInterval
	}

	return config
}

func parseAddress(config *Config) {
	args := strings.Split(config.Host, ":")
	if len(args) == 2 {
		if len(args[0]) == 0 || len(args[1]) == 0 {
			config.Host = host + ":" + port
		}
	} else {
		config.Host = host + ":" + port
	}
}
