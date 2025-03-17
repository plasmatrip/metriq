package db

import (
	"context"
	"embed"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/plasmatrip/metriq/internal/logger"
	"github.com/plasmatrip/metriq/internal/models"
	"github.com/plasmatrip/metriq/internal/types"
)

type PostgresStorage struct {
	db *pgxpool.Pool
	lg logger.Logger
}

func NewPostgresStorage(ctx context.Context, dsn string, lg logger.Logger) (*PostgresStorage, error) {
	// запускаем миграцию
	err := startMigration(dsn)
	if err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			return nil, err
		} else {
			lg.Sugar.Debugw("the database exists, there is nothing to migrate")
		}
	} else {
		lg.Sugar.Debugw("database migration was successful")
	}

	// открываем БД
	db, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, err
	}

	ps := &PostgresStorage{
		db: db,
		lg: lg,
	}

	// // создаем таблицу, при ошибке прокидываем ее наверх
	// err = ps.createTables(ctx)
	// if err != nil {
	// 	return nil, err
	// }

	return ps, nil
}

//go:embed migrations/*.sql
var migrationsDir embed.FS

// StartMigration запускает миграцию
func startMigration(dsn string) error {
	d, err := iofs.New(migrationsDir, "migrations")
	if err != nil {
		return fmt.Errorf("failed to return an iofs driver: %w", err)
	}

	m, err := migrate.NewWithSourceInstance("iofs", d, dsn)
	if err != nil {
		return fmt.Errorf("failed to get a new migrate instance: %w", err)
	}
	if err := m.Up(); err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			return fmt.Errorf("failed to apply migrations to the DB: %w", err)
		}
		return err
	}
	return nil
}

// func (ps PostgresStorage) createTables(ctx context.Context) error {
// 	_, err := ps.db.Exec(ctx, schema)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

func (ps PostgresStorage) Ping(ctx context.Context) error {
	return ps.db.Ping(ctx)
}

func (ps PostgresStorage) Close() {
	ps.db.Close()
}

func (ps PostgresStorage) SetMetrics(ctx context.Context, metrics []models.Metrics) error {
	// начинаем транзакцию
	tx, err := ps.db.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted, AccessMode: pgx.ReadWrite})
	if err != nil {
		return err
	}

	// при ошибке коммита откатываем назад
	defer func() error {
		return tx.Rollback(ctx)
	}()

	// итерируемся по метрикам
	for _, metric := range metrics {
		switch metric.MType {
		case types.Gauge:
			// пытаемся обновить метрику в БД, при ошибке прокидываем ее наверх
			_, err = tx.Exec(ctx, insertGauge, pgx.NamedArgs{
				"id":    metric.ID,
				"mType": metric.MType,
				"value": metric.Value,
			})
			if err != nil {
				return err
			}
			// т.к. пришел тип gauge, увеличиваем PollCounter на 1
			_, err = tx.Exec(ctx, insertCounter, pgx.NamedArgs{
				"id":    types.PollCount,
				"mType": types.Counter,
				"delta": 1,
			})
			if err != nil {
				return err
			}
		case types.Counter:
			_, err = tx.Exec(ctx, insertCounter, pgx.NamedArgs{
				"id":    metric.ID,
				"mType": metric.MType,
				"delta": metric.Delta,
			})
			if err != nil {
				return err
			}
		}
	}

	// запускаем коммит
	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (ps PostgresStorage) SetMetric(ctx context.Context, id string, metric types.Metric) error {
	// определяем тип пришедшей метрики
	switch metric.MetricType {
	case types.Gauge:
		// проверяем метрику (тип и значение)
		if err := metric.Check(); err != nil {
			return err
		}

		// пытаемся обновить метрику в БД, при ошибке прокидываем ее наверх
		res, err := ps.db.Exec(ctx, insertGauge,
			pgx.NamedArgs{
				"id":    id,
				"mType": metric.MetricType,
				"value": metric.Value,
			})
		if err != nil {
			return err
		}
		rows := res.RowsAffected()
		if rows == 0 {
			return errors.New("zero rows inserted")
		}

		// т.к. пришел тип gauge, увеличиваем PollCounter на 1
		err = ps.setCounter(ctx, types.PollCount, types.Metric{MetricType: types.Counter, Value: 1})
		if err != nil {
			return err
		}
	case types.Counter:
		err := ps.setCounter(ctx, id, metric)
		if err != nil {
			return err
		}
	}

	return nil
}

func (ps PostgresStorage) setCounter(ctx context.Context, id string, metric types.Metric) error {
	// пытаемся обновить метрику в БД, при ошибке прокидываем ее наверх
	res, err := ps.db.Exec(ctx, insertCounter,
		pgx.NamedArgs{
			"id":    id,
			"mType": metric.MetricType,
			"delta": metric.Value,
		})
	if err != nil {
		return err
	}

	rows := res.RowsAffected()
	if rows == 0 {
		return errors.New("zero rows inserted")
	}

	return nil
}

func (ps PostgresStorage) Metric(ctx context.Context, id string) (types.Metric, error) {
	m := models.Metrics{}

	// делаем запрос в БД
	row := ps.db.QueryRow(ctx, "SELECT * FROM metrics WHERE id = @id", pgx.NamedArgs{"id": id})

	// читаем результат в структуру models.Metrics, при ошибке прокидываем ее наверх
	err := row.Scan(&m.ID, &m.MType, &m.Value, &m.Delta)
	if err != nil {
		return types.Metric{}, err
	}

	// определяем какая метрика получена и заполняем структуру types.Metric для ответа
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

func (ps PostgresStorage) Metrics(ctx context.Context) (map[string]types.Metric, error) {
	// создаем мапу для записи результата
	metrics := make(map[string]types.Metric, 0)

	// делаем запрос в БД
	rows, err := ps.db.Query(ctx, "SELECT * FROM metrics")

	// при ошибке прокидываем ее наверх
	if err != nil {
		return metrics, err
	}
	defer rows.Close()

	// итерируемся по строкам
	for rows.Next() {
		m := models.Metrics{}
		// читаем результат в структуру models.Metrics, при ошибке прокидываем ее наверх
		err := rows.Scan(&m.ID, &m.MType, &m.Value, &m.Delta)
		if err != nil {
			return metrics, err
		}

		// определяем какая метрика получена и заполняем структуру types.Metric для ответа
		metric := types.Metric{}
		switch m.MType {
		case types.Gauge:
			metric.MetricType = types.Gauge
			metric.Value = *m.Value
		case types.Counter:
			metric.MetricType = types.Counter
			metric.Value = *m.Delta
		}

		// добавляем метрику в мапу
		metrics[m.ID] = metric
	}

	// если при итерации по строкам была ошибка
	if err := rows.Err(); err != nil {
		return metrics, err
	}

	return metrics, nil
}

func (ps PostgresStorage) SetBackup(chan struct{}) {

}
