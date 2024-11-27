package storage

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/plasmatrip/metriq/internal/types"
)

type PosrgresStorage struct {
	DB *sql.DB
}

func NewPostgresStorage(dsn string) (*PosrgresStorage, error) {
	// dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
	// 	`localhost`, `metriq`, ``, `metriq`)

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	return &PosrgresStorage{
		DB: db,
	}, nil
}

func (ps PosrgresStorage) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := ps.DB.PingContext(ctx); err != nil {
		return err
	}
	return nil
}

func (ps PosrgresStorage) Close() error {
	err := ps.DB.Close()
	if err != nil {
		return err
	}
	return nil
}

func (ps PosrgresStorage) SetMetric(key string, metric types.Metric) error {
	return nil
}

func (ps PosrgresStorage) Metric(key string) (types.Metric, bool) {
	return types.Metric{}, false
}

func (ps PosrgresStorage) Metrics() map[string]types.Metric {
	storage := make(map[string]types.Metric, 0)
	return storage
}

func (ps PosrgresStorage) SetBackup(chan struct{}) {

}
