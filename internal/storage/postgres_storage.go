package storage

import (
	"context"
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/plasmatrip/metriq/internal/types"
)

type PosrgresStorage struct {
	db *sql.DB
}

func NewPostgresStorage(dsn string) (*PosrgresStorage, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	return &PosrgresStorage{
		db: db,
	}, nil
}

func (ps PosrgresStorage) Ping(ctx context.Context) error {
	//ctx, cancel := context.WithCancel(context.Background()) //, 1*time.Second)
	//defer cancel()
	// if err := ps.DB.PingContext(ctx); err != nil {
	// 	return err
	// }
	// return nil
	return ps.db.PingContext(ctx)
}

func (ps PosrgresStorage) Close() error {
	err := ps.db.Close()
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
