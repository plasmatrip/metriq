package db

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/plasmatrip/metriq/internal/models"
	"github.com/plasmatrip/metriq/internal/types"
)

type PostgresStorage struct {
	db *sql.DB
}

func NewPostgresStorage(dsn string) (*PostgresStorage, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	ps := &PostgresStorage{
		db: db,
	}

	err = ps.createTables()
	if err != nil {
		return nil, err
	}

	return ps, nil
}

func (ps PostgresStorage) createTables() error {
	_, err := ps.db.Exec(schema)
	if err != nil {
		return err
	}
	return nil
}

func (ps PostgresStorage) Ping(ctx context.Context) error {
	return ps.db.PingContext(ctx)
}

func (ps PostgresStorage) Close() error {
	err := ps.db.Close()
	if err != nil {
		return err
	}
	return nil
}

func (ps PostgresStorage) SetMetric(id string, metric types.Metric) error {
	switch metric.MetricType {
	case types.Gauge:
		if err := metric.Check(); err != nil {
			return err
		}

		res, err := ps.db.Exec("UPDATE metrics SET value = @value WHERE id = @id",
			pgx.NamedArgs{
				"id":    id,
				"value": metric.Value,
			})
		if err != nil {
			return err
		}

		rows, _ := res.RowsAffected()
		if rows == 0 {
			res, err := ps.db.Exec("INSERT INTO metrics (id, mType, value) VALUES (@id, @mType, @value)",
				pgx.NamedArgs{
					"id":    id,
					"mType": metric.MetricType,
					"value": metric.Value,
				})
			if err != nil {
				return err
			}
			rows, _ := res.RowsAffected()
			if rows == 0 {
				return errors.New("zero rows inserted")
			}
		}

		err = ps.setCounter(types.PollCount, types.Metric{MetricType: types.Counter, Value: int64(1)})
		if err != nil {
			return err
		}
	case types.Counter:
		err := ps.setCounter(id, metric)
		if err != nil {
			return err
		}
	}

	return nil
}

func (ps PostgresStorage) setCounter(id string, metric types.Metric) error {
	res, err := ps.db.Exec("UPDATE metrics SET delta = delta + @delta WHERE id = @id",
		pgx.NamedArgs{
			"id":    id,
			"delta": metric.Value,
		})
	if err != nil {
		return err
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		res, err := ps.db.Exec("INSERT INTO metrics (id, mType, delta) VALUES (@id, @mType, @delta)",
			pgx.NamedArgs{
				"id":    id,
				"mType": metric.MetricType,
				"delta": 1,
			})
		if err != nil {
			return err
		}

		rows, _ := res.RowsAffected()
		if rows == 0 {
			return errors.New("zero rows inserted")
		}
	}

	return nil
}

func (ps PostgresStorage) Metric(id string) (types.Metric, error) {
	m := models.Metrics{}

	row := ps.db.QueryRow("SELECT * FROM metrics WHERE id = @id", pgx.NamedArgs{"id": id})

	err := row.Scan(&m.ID, &m.MType, &m.Value, &m.Delta)
	if err != nil {
		return types.Metric{}, err
	}

	if err := row.Err(); err != nil {
		return types.Metric{}, err
	}

	metric := types.Metric{}
	switch m.MType {
	case types.Gauge:
		metric.MetricType = types.Gauge
		metric.Value = *m.Value
	case types.Counter:
		metric.MetricType = types.Counter
		metric.Value = *m.Delta
	}
	return metric, nil
}

func (ps PostgresStorage) Metrics() (map[string]types.Metric, error) {
	metrics := make(map[string]types.Metric, 0)

	rows, err := ps.db.Query("SELECT * FROM metrics")

	if err != nil {
		return metrics, err
	}
	defer rows.Close()

	for rows.Next() {
		m := models.Metrics{}
		err := rows.Scan(&m.ID, &m.MType, &m.Value, &m.Delta)
		if err != nil {
			return metrics, err
		}

		metric := types.Metric{}
		switch m.MType {
		case types.Gauge:
			metric.MetricType = types.Gauge
			metric.Value = *m.Value
		case types.Counter:
			metric.MetricType = types.Counter
			metric.Value = *m.Delta
		}

		metrics[m.ID] = metric
	}

	if err := rows.Err(); err != nil {
		return metrics, err
	}

	return metrics, nil
}

func (ps PostgresStorage) SetBackup(chan struct{}) {

}
