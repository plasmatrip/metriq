package server

const (
	Port    = "8080"
	Address = "localhost"

	Url = "http://" + Address + ":" + Port

	updateURILen = 5
	mTypePos     = 2
	mNamePos     = 3
	mValuePos    = 4

	Gauge   = "gauge"
	Counter = "counter"
)
