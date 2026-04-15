CREATE INDEX idx_reservd_court ON reserv_d(court_uuid);
CREATE INDEX idx_reservd_time ON reserv_d(start_reserv_date, end_reserv_date);
CREATE INDEX idx_payment_reserv ON payment(reserv_h_uuid);