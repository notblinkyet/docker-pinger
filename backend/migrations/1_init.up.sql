CREATE TABLE IF NOT EXISTS containers(
    id SERIAL PRIMARY KEY,
    ip VARCHAR(15) UNIQUE,
    is_tracked BOOLEAN
);

CREATE TABLE IF NOT EXISTS pings(
    id SERIAL PRIMARY KEY,
    container_id INT,
    latency BIGINT NOT NULL,
    last_success_at TIMESTAMP DEFAULT NULL,
    ping_at TIMESTAMP,
    FOREIGN KEY (container_id) REFERENCES containers (id)
);

;
CREATE TABLE IF NOT EXISTS last_pings(
    id INT PRIMARY KEY,
    container_id INT UNIQUE,
    latency BIGINT NOT NULL,
    last_success_at TIMESTAMP DEFAULT NULL,
    ping_at TIMESTAMP,
    FOREIGN KEY (container_id) REFERENCES containers (id) ON DELETE CASCADE
);