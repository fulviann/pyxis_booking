-- ======================
-- DROP TABLES (ORDER: CHILD → PARENT)
-- ======================

DROP TABLE IF EXISTS webhook_logs;
DROP TABLE IF EXISTS payment;
DROP TABLE IF EXISTS reserv_d;
DROP TABLE IF EXISTS reserv_h;
DROP TABLE IF EXISTS cart;
DROP TABLE IF EXISTS blackouts;
DROP TABLE IF EXISTS court;
DROP TABLE IF EXISTS member;
DROP TABLE IF EXISTS merchants;
DROP TABLE IF EXISTS master_sequence;

-- ======================
-- DROP ENUM
-- ======================
DROP TYPE IF EXISTS webhook_status;

-- ======================
-- OPTIONAL: DROP EXTENSION
-- ======================
-- Hapus hanya jika yakin tidak dipakai di tempat lain
-- DROP EXTENSION IF EXISTS "pgcrypto";