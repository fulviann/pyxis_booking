-- ======================
-- EXTENSION
-- ======================
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- ======================
-- ENUM (PostgreSQL)
-- ======================
CREATE TYPE webhook_status AS ENUM ('pending','sent','failed');

-- ======================
-- TABLE: merchants
-- ======================
CREATE TABLE merchants (
    merchant_uuid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    merchant_code VARCHAR(13) UNIQUE NOT NULL,
    name VARCHAR(100) NOT NULL,
    address TEXT,
    phone_number VARCHAR(15),
    email VARCHAR(255),
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ======================
-- TABLE: court
-- ======================
CREATE TABLE court (
    court_uuid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    court_code VARCHAR(13) UNIQUE NOT NULL,
    merchant_uuid UUID NOT NULL,
    name VARCHAR(100) NOT NULL,
    type VARCHAR(100) NOT NULL,
    price_per_hour DECIMAL(10,2) NOT NULL,
    status VARCHAR(50) NOT NULL,
    capacity INT NOT NULL DEFAULT 10,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_court_merchant
        FOREIGN KEY (merchant_uuid)
        REFERENCES merchants(merchant_uuid)
        ON DELETE CASCADE
);

-- ======================
-- TABLE: member
-- ======================
CREATE TABLE member (
    member_uuid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    member_code VARCHAR(13) UNIQUE NOT NULL,
    name VARCHAR(100) NOT NULL,
    phone_number VARCHAR(15) NOT NULL,
    email VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    position VARCHAR(100) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ======================
-- TABLE: cart
-- ======================
CREATE TABLE cart (
    cart_uuid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    cart_code VARCHAR(13) UNIQUE NOT NULL,
    court_uuid UUID,
    member_uuid UUID,
    duration INT NOT NULL,
    price DECIMAL(10,2) NOT NULL,
    start_reserv_date TIMESTAMPTZ,
    end_reserv_date TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_cart_court
        FOREIGN KEY (court_uuid)
        REFERENCES court(court_uuid)
        ON DELETE SET NULL,

    CONSTRAINT fk_cart_member
        FOREIGN KEY (member_uuid)
        REFERENCES member(member_uuid)
        ON DELETE SET NULL
);

-- ======================
-- TABLE: reserv_h
-- ======================
CREATE TABLE reserv_h (
    reserv_h_uuid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    member_uuid UUID,
    reserv_h_code VARCHAR(13) UNIQUE NOT NULL,

    member_name VARCHAR(100) NOT NULL,
    phone_number VARCHAR(15) NOT NULL,
    email VARCHAR(255),

    player_qty INT NOT NULL,
    total_amount DECIMAL(10,2) NOT NULL,
    status VARCHAR(50) NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_reserv_member
        FOREIGN KEY (member_uuid)
        REFERENCES member(member_uuid)
        ON DELETE SET NULL
);

-- ======================
-- TABLE: reserv_d
-- ======================
CREATE TABLE reserv_d (
    reserv_d_uuid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    reserv_h_uuid UUID,
    court_uuid UUID,

    price_per_hour DECIMAL(10,2) NOT NULL,
    total_amount DECIMAL(10,2) NOT NULL,
    duration INT NOT NULL,
    start_reserv_date TIMESTAMPTZ NOT NULL,
    end_reserv_date TIMESTAMPTZ NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_reservd_reservh
        FOREIGN KEY (reserv_h_uuid)
        REFERENCES reserv_h(reserv_h_uuid)
        ON DELETE CASCADE,

    CONSTRAINT fk_reservd_court
        FOREIGN KEY (court_uuid)
        REFERENCES court(court_uuid),

    CONSTRAINT chk_time_valid
        CHECK (end_reserv_date > start_reserv_date)
);

-- ======================
-- TABLE: payment
-- ======================
CREATE TABLE payment (
    payment_uuid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    payment_code VARCHAR(13) UNIQUE,
    reserv_h_uuid UUID NOT NULL,

    method VARCHAR(100) NOT NULL,
    status VARCHAR(50) NOT NULL,
    total_amount DECIMAL(10,2) NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_payment_reservh
        FOREIGN KEY (reserv_h_uuid)
        REFERENCES reserv_h(reserv_h_uuid)
        ON DELETE CASCADE
);

-- ======================
-- TABLE: master_sequence
-- ======================
CREATE TABLE master_sequence (
    master_sequence_uuid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    master_seq_code VARCHAR(13) UNIQUE NOT NULL,
    sequence_name VARCHAR(100) NOT NULL,
    seq_no VARCHAR(15) NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ======================
-- TABLE: webhook_logs
-- ======================
CREATE TABLE webhook_logs (
    log_uuid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    reserv_h_uuid UUID NOT NULL,
    event_type VARCHAR(50) NOT NULL,
    target_type VARCHAR(20) NOT NULL,
    target_addr VARCHAR(255) NOT NULL,
    payload JSON,
    status webhook_status NOT NULL DEFAULT 'pending',
    sent_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_webhook_logs_reservh
        FOREIGN KEY (reserv_h_uuid)
        REFERENCES reserv_h(reserv_h_uuid)
        ON DELETE CASCADE
);

-- ======================
-- TABLE: blackouts
-- ======================
CREATE TABLE blackouts (
    blackout_uuid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    court_uuid UUID,
    merchant_uuid UUID,
    title VARCHAR(100) NOT NULL,
    start_date TIMESTAMPTZ NOT NULL,
    end_date TIMESTAMPTZ NOT NULL,
    is_full_day BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_blackouts_court
        FOREIGN KEY (court_uuid)
        REFERENCES court(court_uuid)
        ON DELETE CASCADE,

    CONSTRAINT fk_blackouts_merchant
        FOREIGN KEY (merchant_uuid)
        REFERENCES merchants(merchant_uuid)
        ON DELETE CASCADE,

    CONSTRAINT chk_blackout_time
        CHECK (end_date > start_date)
);