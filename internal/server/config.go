package server

import (
	"flag"
	"fmt"
	"strings"
)

const (
	Port = "8080"
	Host = "localhost"

	URL = "http://" + Host + ":" + Port

	UpdateURILen = 5
	ValueURILen  = 4

	RequestTypePos  = 2
	RequestNamePos  = 3
	RequestValuePos = 4

	Gauge   = "gauge"
	Counter = "counter"

	PollCount = "PollCount"
)

type Config struct {
	Port string
	Host string
}

func NewConfig() *Config {
	var srv string
	flag.StringVar(&srv, "a", "localhost:8080", "Server address host:port")
	flag.Parse()
	args := strings.Split(srv, ":")

	fmt.Println(srv)

	fmt.Println(args)

	// server := new(types.SrvAddr)
	// _ = flag.Value(server)
	// flag.Var(server, "a", "Server address host:port")
	// flag.Parse()

	// if len(flag.Args()) > 0 {
	// 	fmt.Println("Unknown flag(s): ", flag.Args())
	// 	os.Exit(1)
	// }

	// if len(server.Host) == 0 || len(server.Port) == 0 {
	// 	server.Host = Host
	// 	server.Port = Port
	// }

	// if len(args)

	return &Config{
		Host: args[0],
		Port: args[1],
	}
}
