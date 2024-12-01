package db

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/plasmatrip/metriq/internal/logger"
	"github.com/plasmatrip/metriq/internal/models"
	"github.com/plasmatrip/metriq/internal/types"
)

type PostgresStorage struct {
	db *sql.DB
	lg logger.Logger
}

func NewPostgresStorage(dsn string, lg logger.Logger) (*PostgresStorage, error) {
	// открываем БД
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	ps := &PostgresStorage{
		db: db,
		lg: lg,
	}

	// создаем таблицу, при ошибке прокидываем ее наверх
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
	// проверяем, что есть коннект с БД
	return ps.db.PingContext(ctx)
}

func (ps PostgresStorage) Close() error {
	// закрываем коннект с БД
	err := ps.db.Close()
	if err != nil {
		return err
	}
	return nil
}

func (ps PostgresStorage) SetMetrics(ctx context.Context, metrics models.SMetrics) error {
	// начинаем транзакцию
	tx, err := ps.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	// при ошибке коммита откатываем назад
	defer func() error {
		return tx.Rollback()
	}()

	// итерируемся по метрикам
	for _, metric := range metrics.Metrics {
		switch metric.MType {
		case types.Gauge:
			// пытаемся обновить метрику в БД, при ошибке прокидываем ее наверх
			_, err := tx.ExecContext(ctx, insertGauge, pgx.NamedArgs{
				"id":    metric.ID,
				"mType": metric.MType,
				"value": metric.Value,
			})
			if err != nil {
				return err
			}
			// т.к. пришел тип gauge, увеличиваем PollCounter на 1
			_, err = tx.ExecContext(ctx, insertCounter, pgx.NamedArgs{
				"id":    types.PollCount,
				"mType": types.Counter,
				"delta": 1,
			})
			if err != nil {
				return err
			}
		case types.Counter:
			_, err := tx.ExecContext(ctx, insertCounter, pgx.NamedArgs{
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
	tx.Commit()
	return nil
}

func (ps PostgresStorage) SetMetric(id string, metric types.Metric) error {
	// определяем тип пришедшей метрики
	switch metric.MetricType {
	case types.Gauge:
		// проверяем метрику (тип и значение)
		if err := metric.Check(); err != nil {
			return err
		}

		// пытаемся обновить метрику в БД, при ошибке прокидываем ее наверх
		res, err := ps.db.Exec(insertGauge,
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

		// т.к. пришел тип gauge, увеличиваем PollCounter на 1
		err = ps.setCounter(types.PollCount, types.Metric{MetricType: types.Counter, Value: 1})
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
	// пытаемся обновить метрику в БД, при ошибке прокидываем ее наверх
	res, err := ps.db.Exec(insertCounter,
		pgx.NamedArgs{
			"id":    id,
			"mType": metric.MetricType,
			"delta": metric.Value,
		})
	if err != nil {
		return err
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return errors.New("zero rows inserted")
	}

	return nil
}

func (ps PostgresStorage) Metric(id string) (types.Metric, error) {
	m := models.Metrics{}

	// делаем запрос в БД
	row := ps.db.QueryRow("SELECT * FROM metrics WHERE id = @id", pgx.NamedArgs{"id": id})

	ps.lg.Sugar.Infoln("select", "row", row)

	// читаем результат в структуру models.Metrics, при ошибке прокидываем ее наверх
	err := row.Scan(&m.ID, &m.MType, &m.Value, &m.Delta)
	if err != nil {
		return types.Metric{}, err
	}

	ps.lg.Sugar.Infoln("select", "metric", m)

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

func (ps PostgresStorage) Metrics() (map[string]types.Metric, error) {
	// создаем мапу для записи результата
	metrics := make(map[string]types.Metric, 0)

	// делаем запрос в БД
	rows, err := ps.db.Query("SELECT * FROM metrics")

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
