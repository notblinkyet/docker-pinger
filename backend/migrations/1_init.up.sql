CREATE TABLE IF NOT EXISTS container(
    id SERIAL PRIMARY KEY,
    ip VARCHAR(15) UNIQUE
);

CREATE TABLE IF NOT EXISTS ping(
    id SERIAL PRIMARY KEY,
    container_id INT,
    latency INT NOT NULL,
    last_success_at TIMESTAMP DEFAULT NULL,
    ping_at TIMESTAMP DEFAULT NOW(),
    FOREIGN KEY (container_id) REFERENCES container (id)
);

CREATE TABLE IF NOT EXISTS last_ping(
    id INT PRIMARY KEY,
    container_id INT UNIQUE,
    latency INT NOT NULL,
    last_success_at TIMESTAMP DEFAULT NULL,
    ping_at TIMESTAMP DEFAULT NOW(),
    FOREIGN KEY (container_id) REFERENCES container (id) ON DELETE CASCADE
);