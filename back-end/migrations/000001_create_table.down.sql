-- ======================
-- DROP TABLES (URUTAN PENTING)
-- ======================
DROP TABLE IF EXISTS payment;
DROP TABLE IF EXISTS reserv_d;
DROP TABLE IF EXISTS reserv_h;
DROP TABLE IF EXISTS court;
DROP TABLE IF EXISTS member;
DROP TABLE IF EXISTS master_sequence;

-- ======================
-- OPTIONAL: DROP EXTENSION
-- ======================
DROP EXTENSION IF EXISTS "pgcrypto";