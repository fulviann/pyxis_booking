-- ======================
-- ENUM (PostgreSQL)
-- ======================
CREATE TYPE webhook_status AS ENUM ('pending', 'sent', 'failed');

-- ======================
-- TABLE: webhook_logs
-- ======================
CREATE TABLE
    webhook_logs (
        log_uuid UUID PRIMARY KEY DEFAULT gen_random_uuid (),
        reserv_h_uuid UUID NOT NULL,
        event_type VARCHAR(50) NOT NULL,
        target_type VARCHAR(20) NOT NULL,
        target_addr VARCHAR(255) NOT NULL,
        payload JSON,
        status webhook_status NOT NULL DEFAULT 'pending',
        sent_at TIMESTAMPTZ,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        CONSTRAINT fk_webhook_logs_reservh FOREIGN KEY (reserv_h_uuid) REFERENCES reserv_h (reserv_h_uuid) ON DELETE CASCADE
    );
