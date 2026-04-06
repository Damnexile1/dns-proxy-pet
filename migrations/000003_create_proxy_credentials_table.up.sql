-- Create proxy_credentials table
CREATE TABLE IF NOT EXISTS proxy_credentials (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token VARCHAR(255) UNIQUE NOT NULL,
    login VARCHAR(255) UNIQUE,
    password VARCHAR(255),
    max_devices INTEGER DEFAULT 1,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_proxy_credentials_user_id ON proxy_credentials(user_id);
CREATE INDEX idx_proxy_credentials_token ON proxy_credentials(token);
CREATE INDEX idx_proxy_credentials_login ON proxy_credentials(login);
