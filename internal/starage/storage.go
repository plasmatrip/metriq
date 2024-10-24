package starage

type MemStorage struct {
	Gouge   map[string]float64
	Counter map[string]int64
}

type Repository interface {
	Add(key string, value any)
	Delete(key string)
	Update(key string, value int64)
}
