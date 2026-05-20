-- ======================
-- DROP CHILD TABLES FIRST
-- ======================

DROP TABLE IF EXISTS webhook_logs;
DROP TABLE IF EXISTS invalid_token;

DROP TABLE IF EXISTS reserv_payment;
DROP TABLE IF EXISTS reserv_d;
DROP TABLE IF EXISTS cart;
DROP TABLE IF EXISTS blackouts;

DROP TABLE IF EXISTS services;

DROP TABLE IF EXISTS merchants;

DROP TABLE IF EXISTS subscription_payment;

DROP TABLE IF EXISTS reserv_h;

DROP TABLE IF EXISTS transactions;

DROP TABLE IF EXISTS subscription_list;

-- ======================
-- DROP MASTER TABLES
-- ======================

DROP TABLE IF EXISTS customer;
DROP TABLE IF EXISTS admin;
DROP TABLE IF EXISTS owner;

DROP TABLE IF EXISTS packages;
DROP TABLE IF EXISTS service_category;

DROP TABLE IF EXISTS master_sequence;

-- ======================
-- DROP ENUM
-- ======================

DROP TYPE IF EXISTS webhook_status;