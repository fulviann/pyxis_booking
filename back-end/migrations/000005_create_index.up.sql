-- ======================
-- INDEXES
-- ======================
CREATE INDEX idx_reservd_service ON reserv_d(service_uuid);
CREATE INDEX idx_reservd_time ON reserv_d(start_reserv_date, end_reserv_date);
CREATE INDEX idx_payment_reserv ON reserv_payment(reserv_h_uuid);
CREATE INDEX idx_blackout_range ON blackouts(service_uuid, start_date, end_date);

-- ======================
-- EXTENSIONS
-- ======================
CREATE EXTENSION IF NOT EXISTS btree_gist;

ALTER TABLE reserv_d
ADD CONSTRAINT no_double_booking
EXCLUDE USING gist (
    service_uuid WITH =,
    tstzrange(start_reserv_date, end_reserv_date) WITH &&
);

-- ======================
-- EXTENSION
-- ======================
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- ======================
-- TABLE: master_sequence
-- ======================
CREATE TABLE
    master_sequence (
        master_sequence_uuid UUID PRIMARY KEY DEFAULT gen_random_uuid (),
        master_seq_code VARCHAR(15) UNIQUE NOT NULL,
        sequence_name VARCHAR(100) NOT NULL,
        seq_no VARCHAR(15) NOT NULL,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW ()
    );