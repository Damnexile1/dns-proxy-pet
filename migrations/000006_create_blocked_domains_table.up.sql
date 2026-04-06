-- Create blocked_domains table
CREATE TABLE IF NOT EXISTS blocked_domains (
    id SERIAL PRIMARY KEY,
    domain VARCHAR(255) UNIQUE NOT NULL,
    added_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_blocked_domains_domain ON blocked_domains(domain);
