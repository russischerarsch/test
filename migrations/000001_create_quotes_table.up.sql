CREATE TABLE IF NOT EXISTS quotes(
		id SERIAL PRIMARY KEY,
		pair VARCHAR(7) NOT NULL UNIQUE,
		price DECIMAL,
        created_at TIMESTAMP NOT NULL,
		updated_at TIMESTAMP,
		status VARCHAR(10) DEFAULT 'pending'
	)