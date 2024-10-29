package storage

type Repository interface {
	Update(key string, metric Metric) error
	Get(key string) (Metric, bool)
	GetAll() map[string]Metric
}
