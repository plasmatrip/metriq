package agent

import (
	"flag"
	"fmt"
	"os"
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

	var srv string
	flag.StringVar(&srv, "a", "localhost:8080", "Server address host:port")
	//flag.Parse()
	//args := strings.Split(srv, ":")

	//fmt.Println(srv)

	//fmt.Println(args)

	// server := new(types.SrvAddr)
	// _ = flag.Value(server)
	// flag.Var(server, "a", "Server address host:port")
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

	// if len(server.Host) == 0 || len(server.Port) == 0 {
	// 	server.Host = host
	// 	server.Port = port
	// }

	//config.Host = args[0]
	//config.Port = args[1]
	config.URL = "http://" + srv //"http://" + config.Host + ":" + config.Port

	return config
}
