-- Create traffic_usage table
CREATE TABLE IF NOT EXISTS traffic_usage (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    bytes_used BIGINT DEFAULT 0,
    date DATE NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_traffic_usage_user_id ON traffic_usage(user_id);
CREATE INDEX idx_traffic_usage_date ON traffic_usage(date);
CREATE UNIQUE INDEX idx_traffic_usage_user_date ON traffic_usage(user_id, date);
