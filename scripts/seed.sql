-- Seed file for test data

-- Insert test users
INSERT INTO users (telegram_id, username, created_at, updated_at) VALUES
(123456789, 'test_user_1', NOW(), NOW()),
(987654321, 'test_user_2', NOW(), NOW()),
(555555555, 'test_user_3', NOW(), NOW());

-- Insert test subscriptions
INSERT INTO subscriptions (user_id, plan_type, status, started_at, expires_at, auto_renew) VALUES
(1, 'basic', 'active', NOW(), NOW() + INTERVAL '30 days', TRUE),
(2, 'premium', 'active', NOW(), NOW() + INTERVAL '30 days', TRUE),
(3, 'family', 'expired', NOW() - INTERVAL '60 days', NOW() - INTERVAL '30 days', FALSE);

-- Insert test proxy credentials
INSERT INTO proxy_credentials (user_id, token, login, password, max_devices) VALUES
(1, 'test_token_user1_abc123', 'user1', 'password1', 1),
(2, 'test_token_user2_def456', 'user2', 'password2', 3),
(3, 'test_token_user3_ghi789', 'user3', 'password3', 5);

-- Insert test payments
INSERT INTO payments (user_id, subscription_id, amount, currency, payment_method, status, external_id) VALUES
(1, 1, 300.00, 'RUB', 'yookassa', 'completed', 'yookassa_payment_1'),
(2, 2, 500.00, 'RUB', 'yookassa', 'completed', 'yookassa_payment_2'),
(3, 3, 800.00, 'RUB', 'yookassa', 'completed', 'yookassa_payment_3');

-- Insert test traffic usage
INSERT INTO traffic_usage (user_id, bytes_used, date) VALUES
(1, 1073741824, CURRENT_DATE),  -- 1 GB
(2, 5368709120, CURRENT_DATE),  -- 5 GB
(3, 10737418240, CURRENT_DATE); -- 10 GB

-- Insert blocked domains
INSERT INTO blocked_domains (domain) VALUES
('facebook.com'),
('*.facebook.com'),
('instagram.com'),
('*.instagram.com'),
('twitter.com'),
('*.twitter.com'),
('youtube.com'),
('*.youtube.com'),
('linkedin.com'),
('*.linkedin.com');

-- Display inserted data
SELECT 'Users:' as info;
SELECT id, telegram_id, username FROM users;

SELECT 'Subscriptions:' as info;
SELECT id, user_id, plan_type, status, expires_at FROM subscriptions;

SELECT 'Proxy Credentials:' as info;
SELECT id, user_id, login, max_devices FROM proxy_credentials;

SELECT 'Payments:' as info;
SELECT id, user_id, amount, currency, status FROM payments;

SELECT 'Traffic Usage:' as info;
SELECT id, user_id, bytes_used, date FROM traffic_usage;

SELECT 'Blocked Domains:' as info;
SELECT id, domain FROM blocked_domains LIMIT 10;
