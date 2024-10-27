package types

import (
	"errors"
	"strings"
)

const (
	Gauge   = "gauge"
	Counter = "counter"

	PollCount = "PollCount"
)

type SrvAddr struct {
	Host string
	Port string
}

func (srv SrvAddr) String() string {
	return srv.Host + ":" + srv.Port
}

func (srv *SrvAddr) Set(s string) error {
	addr := strings.Split(s, ":")

	if len(addr) != 2 {
		return errors.New("need address in a form host:port")
	}

	srv.Host = addr[0]
	srv.Port = addr[1]

	return nil
}
