package db

var schema = `
		CREATE TABLE IF NOT EXISTS metrics (
			id VARCHAR(64) NOT NULL PRIMARY KEY,
			mType VARCHAR(64) NOT NULL,
			value DOUBLE PRECISION DEFAULT NULL,
			delta INTEGER DEFAULT NULL
		);
		CREATE INDEX IF NOT EXISTS metrics ON metrics (id);
	`
