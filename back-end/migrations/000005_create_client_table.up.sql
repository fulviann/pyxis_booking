-- ======================
-- TABLE: merchants
-- ======================
CREATE TABLE
    merchants (
        merchant_uuid UUID PRIMARY KEY DEFAULT gen_random_uuid (),
        merchant_code VARCHAR(15) UNIQUE NOT NULL,
        subscription_list_uuid UUID NOT NULL,
        name VARCHAR(100) NOT NULL,
        address TEXT,
        phone_number VARCHAR(15),
        email VARCHAR(255),
        is_active BOOLEAN NOT NULL DEFAULT TRUE,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        CONSTRAINT fk_merchant_subscription_list FOREIGN KEY (subscription_list_uuid) REFERENCES subscription_list (subscription_list_uuid) ON DELETE CASCADE
    );

-- ======================
-- TABLE: services
-- ======================
CREATE TABLE
    services (
        service_uuid UUID PRIMARY KEY DEFAULT gen_random_uuid (),
        service_code VARCHAR(15) UNIQUE NOT NULL,
        subscription_list_uuid UUID NOT NULL,
        service_type_uuid UUID NOT NULL,
        merchant_uuid UUID NOT NULL,
        name VARCHAR(100) NOT NULL,
        type VARCHAR(100) NOT NULL,
        price DECIMAL(10, 2) NOT NULL,
        price_type VARCHAR(50) NOT NULL, -- e.g., "minute", "hour", "day"
        description TEXT,
        status VARCHAR(50) NOT NULL,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        CONSTRAINT fk_service_service_type FOREIGN KEY (service_type_uuid) REFERENCES service_types (service_type_uuid) ON DELETE CASCADE,
        CONSTRAINT fk_service_subscription_list FOREIGN KEY (subscription_list_uuid) REFERENCES subscription_list (subscription_list_uuid) ON DELETE CASCADE,
        CONSTRAINT fk_service_merchant FOREIGN KEY (merchant_uuid) REFERENCES merchants (merchant_uuid) ON DELETE CASCADE
    );

-- ======================
-- TABLE: cart
-- ======================
CREATE TABLE
    cart (
        cart_uuid UUID PRIMARY KEY DEFAULT gen_random_uuid (),
        cart_code VARCHAR(15) UNIQUE NOT NULL,
        service_uuid UUID,
        customer_uuid UUID,
        duration INT NOT NULL,
        price DECIMAL(10, 2) NOT NULL,
        price_type VARCHAR(50) NOT NULL, -- e.g., "minute", "hour", "day"
         total_amount DECIMAL(10, 2) NOT NULL,
        start_reserv_date TIMESTAMPTZ,
        end_reserv_date TIMESTAMPTZ,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        CONSTRAINT fk_cart_service FOREIGN KEY (service_uuid) REFERENCES services (service_uuid) ON DELETE SET NULL,
        CONSTRAINT fk_cart_customer FOREIGN KEY (customer_uuid) REFERENCES customer (customer_uuid) ON DELETE SET NULL
    );

-- ======================
-- TABLE: reserv_h
-- ======================
CREATE TABLE
    reserv_h (
        reserv_h_uuid UUID PRIMARY KEY DEFAULT gen_random_uuid (),
        customer_uuid UUID,
        reserv_h_code VARCHAR(15) UNIQUE NOT NULL,
        customer_name VARCHAR(100) NOT NULL,
        phone_number VARCHAR(15) NOT NULL,
        email VARCHAR(255),
        player_qty INT NOT NULL,
        total_amount DECIMAL(10, 2) NOT NULL,
        status VARCHAR(50) NOT NULL,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        CONSTRAINT fk_reserv_customer FOREIGN KEY (customer_uuid) REFERENCES customer (customer_uuid) ON DELETE SET NULL
    );

-- ======================
-- TABLE: reserv_d
-- ======================
CREATE TABLE
    reserv_d (
        reserv_d_uuid UUID PRIMARY KEY DEFAULT gen_random_uuid (),
        reserv_h_uuid UUID,
        service_uuid UUID,
        price DECIMAL(10, 2) NOT NULL,
        price_type VARCHAR(50) NOT NULL -- e.g., "minute", "hour", "day"
        duration INT NOT NULL,
        total_amount DECIMAL(10, 2) NOT NULL,
        start_reserv_date TIMESTAMPTZ NOT NULL,
        end_reserv_date TIMESTAMPTZ NOT NULL,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        CONSTRAINT fk_reservd_reservh FOREIGN KEY (reserv_h_uuid) REFERENCES reserv_h (reserv_h_uuid) ON DELETE CASCADE,
        CONSTRAINT fk_reservd_service FOREIGN KEY (service_uuid) REFERENCES services (service_uuid),
        CONSTRAINT chk_time_valid CHECK (end_reserv_date > start_reserv_date)
    );

-- ======================
-- TABLE: reserv_payment
-- ======================
CREATE TABLE
    reserv_payment (
        payment_uuid UUID PRIMARY KEY DEFAULT gen_random_uuid (),
        payment_code VARCHAR(13) UNIQUE,
        reserv_h_uuid UUID NOT NULL,
        method VARCHAR(100) NOT NULL,
        status VARCHAR(50) NOT NULL,
        total_amount DECIMAL(10, 2) NOT NULL,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        CONSTRAINT fk_payment_reservh FOREIGN KEY (reserv_h_uuid) REFERENCES reserv_h (reserv_h_uuid) ON DELETE CASCADE
    );

-- ======================
-- TABLE: blackouts
-- ======================
CREATE TABLE
    blackouts (
        blackout_uuid UUID PRIMARY KEY DEFAULT gen_random_uuid (),
        service_uuid UUID,
        merchant_uuid UUID,
        title VARCHAR(100) NOT NULL,
        start_date TIMESTAMPTZ NOT NULL,
        end_date TIMESTAMPTZ NOT NULL,
        is_full_day BOOLEAN NOT NULL DEFAULT FALSE,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        CONSTRAINT fk_blackouts_service FOREIGN KEY (service_uuid) REFERENCES services (service_uuid) ON DELETE CASCADE,
        CONSTRAINT fk_blackouts_merchant FOREIGN KEY (merchant_uuid) REFERENCES merchants (merchant_uuid) ON DELETE CASCADE,
        CONSTRAINT chk_blackout_time CHECK (end_date > start_date)
    );