package storage

const (
	PollCount = "PollCount"
)

type (
	Gauge   float64
	Counter int64
)

var Metrics = map[string]struct{}{
	"Alloc":         {},
	"TotalAlloc":    {},
	"Sys":           {},
	"Lookups":       {},
	"Mallocs":       {},
	"Frees":         {},
	"HeapAlloc":     {},
	"HeapSys":       {},
	"HeapIdle":      {},
	"HeapInuse":     {},
	"HeapReleased":  {},
	"HeapObjects":   {},
	"StackInuse":    {},
	"StackSys":      {},
	"MSpanInuse":    {},
	"MSpanSys":      {},
	"MCacheInuse":   {},
	"MCacheSys":     {},
	"BuckHashSys":   {},
	"GCSys":         {},
	"OtherSys":      {},
	"NextGC":        {},
	"LastGC":        {},
	"PauseTotalNs":  {},
	"NumGC":         {},
	"NumForcedGC":   {},
	"GCCPUFraction": {},
	"RandomValue":   {},
	"PollCount":     {},
}

type Repository interface {
	UpdateGauge(key string, value Gauge)
	UpdateCounter(key string, count Counter)
	GetGauges() map[string]Gauge
	GetGauge(key string) Gauge
	GetCounters() map[string]Counter
	GetCounter(key string) Counter
}
