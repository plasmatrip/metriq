package server

import (
	"flag"
	"fmt"
	"os"
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
		fmt.Println("oops")

		return config
	}

	//проверяем флаги
	flag.StringVar(&config.Host, "a", "localhost:8080", "Server address host:port")
	flag.Parse()

	if len(flag.Args()) > 0 {
		fmt.Println("Unknown flag(s): ", flag.Args())
		os.Exit(1)
	}

	parseAddress(config)

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
