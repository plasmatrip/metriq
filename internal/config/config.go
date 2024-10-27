package config

const (
	Port    = "8080"
	Address = "localhost"

	URL = "http://" + Address + ":" + Port

	UpdateURILen = 5
	ValueURILen  = 4

	RequestTypePos  = 2
	RequestNamePos  = 3
	RequestValuePos = 4

	Gauge   = "gauge"
	Counter = "counter"

	PollCount = "PollCount"
)

// type (
// 	Gauge   float64
// 	Counter int64
// )
