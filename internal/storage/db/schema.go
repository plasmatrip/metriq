package db

const (
	schema = `
		CREATE TABLE IF NOT EXISTS metrics (
			id VARCHAR(128) NOT NULL PRIMARY KEY,
			mType VARCHAR(128) NOT NULL,
			value DOUBLE PRECISION DEFAULT NULL,
			delta BIGINT DEFAULT NULL
		);
	`

	insertGauge = `
		INSERT INTO metrics (id, mType, value) VALUES (@id, @mType, @value)
		ON CONFLICT (id)
		DO UPDATE SET value = @value
	`

	insertCounter = `
		INSERT INTO metrics (id, mType, delta) VALUES (@id, @mType, @delta)
		ON CONFLICT (id)
		DO UPDATE SET delta = metrics.delta + @delta
	`
)
