-- ======================
-- INDEXES
-- ======================
CREATE INDEX idx_reservd_court ON reserv_d(court_uuid);
CREATE INDEX idx_reservd_time ON reserv_d(start_reserv_date, end_reserv_date);
CREATE INDEX idx_payment_reserv ON payment(reserv_h_uuid);
CREATE INDEX idx_blackout_range ON blackouts(court_uuid, start_date, end_date);

-- ======================
-- EXTENSIONS
-- ======================
CREATE EXTENSION IF NOT EXISTS btree_gist;

ALTER TABLE reserv_d
ADD CONSTRAINT no_double_booking
EXCLUDE USING gist (
    court_uuid WITH =,
    tstzrange(start_reserv_date, end_reserv_date) WITH &&
);