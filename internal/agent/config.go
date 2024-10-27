package agent

import (
	"flag"
	"fmt"
	"os"

	"github.com/plasmatrip/metriq/internal/types"
)

const (
	port           = "8080"
	host           = "localhost"
	pollInterval   = 2
	reportInterval = 10
)

type Config struct {
	Port           string
	Host           string
	URL            string
	PollInterval   int
	ReportInterval int
}

func NewConfig() *Config {
	config := new(Config)

	server := new(types.SrvAddr)
	_ = flag.Value(server)
	flag.Var(server, "a", "Server address host:port")
	flag.IntVar(&config.PollInterval, "p", pollInterval, "metrics reporting interval")
	flag.IntVar(&config.ReportInterval, "r", reportInterval, "metrics polling frequency")
	flag.Parse()

	if len(flag.Args()) > 0 {
		fmt.Println("Unknown flag(s): ", flag.Args())
		os.Exit(1)
	}

	if config.PollInterval <= 0 {
		config.PollInterval = pollInterval
	}

	if config.ReportInterval <= 0 {
		config.ReportInterval = reportInterval
	}

	if len(server.Host) == 0 || len(server.Port) == 0 {
		server.Host = host
		server.Port = port
	}

	config.Host = server.Host
	config.Port = server.Port
	config.URL = "http://" + server.Host + ":" + server.Port

	return config
}
