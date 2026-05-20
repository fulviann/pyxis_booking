-- ======================
-- EXTENSION
-- ======================
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- ======================
-- TABLE: owner
-- ======================
CREATE TABLE
    owner (
        owner_uuid UUID PRIMARY KEY DEFAULT gen_random_uuid (),
        owner_code VARCHAR(15) UNIQUE NOT NULL,
        name VARCHAR(100) NOT NULL,
        phone_number VARCHAR(15) NOT NULL,
        email VARCHAR(255) NOT NULL,
        password VARCHAR(255) NOT NULL,
        position VARCHAR(100) NOT NULL,
        is_active BOOLEAN NOT NULL DEFAULT TRUE,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW ()
    );

-- ======================
-- TABLE: admin
-- ======================
CREATE TABLE
    admin (
        admin_uuid UUID PRIMARY KEY DEFAULT gen_random_uuid (),
        admin_code VARCHAR(15) UNIQUE NOT NULL,
        name VARCHAR(100) NOT NULL,
        phone_number VARCHAR(15) NOT NULL,
        email VARCHAR(255) NOT NULL,
        password VARCHAR(255) NOT NULL,
        position VARCHAR(100) NOT NULL,
        is_active BOOLEAN NOT NULL DEFAULT TRUE,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW ()
    );

-- ======================
-- TABLE: customer
-- ======================
CREATE TABLE
    customer (
        customer_uuid UUID PRIMARY KEY DEFAULT gen_random_uuid (),
        customer_code VARCHAR(15) UNIQUE NOT NULL,
        name VARCHAR(100) NOT NULL,
        phone_number VARCHAR(15) NOT NULL,
        email VARCHAR(255) NOT NULL,
        password VARCHAR(255) NOT NULL,
        position VARCHAR(100) NOT NULL,
        is_active BOOLEAN NOT NULL DEFAULT TRUE,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW ()
    );

-- ======================
-- TABLE: packages
-- ======================
CREATE TABLE
    packages (
        package_uuid UUID PRIMARY KEY DEFAULT gen_random_uuid (),
        package_code VARCHAR(15) UNIQUE NOT NULL,
        name VARCHAR(100) NOT NULL,
        description TEXT,
        price DECIMAL(10, 2) NOT NULL,
        duration INTEGER NOT NULL,
        is_active BOOLEAN NOT NULL DEFAULT TRUE,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW ()
    );

-- ======================
-- TABLE: service_category
-- ======================
CREATE TABLE
    service_category (
        service_category_uuid UUID PRIMARY KEY DEFAULT gen_random_uuid (),
        service_category_code VARCHAR(15) UNIQUE NOT NULL,
        name VARCHAR(100) NOT NULL, -- e.g., "Hotel", "Lapangan Futsal", "Lapangan Badminton"
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW ()
    );

-- ======================
-- TABLE: subscription_list
-- ======================
CREATE TABLE
    subscription_list (
        subscription_list_uuid UUID PRIMARY KEY DEFAULT gen_random_uuid (),
        package_uuid UUID NOT NULL,
        owner_uuid UUID NOT NULL,
        -- service_category_uuid UUID NOT NULL,
        transaction_date TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        status VARCHAR(50) NOT NULL,
        start_date TIMESTAMPTZ NOT NULL,
        end_date TIMESTAMPTZ NOT NULL,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        CONSTRAINT fk_subscription_list_package FOREIGN KEY (package_uuid) REFERENCES packages (package_uuid) ON DELETE CASCADE,
        CONSTRAINT fk_subscription_list_owner FOREIGN KEY (owner_uuid) REFERENCES owner (owner_uuid) ON DELETE CASCADE,
        -- CONSTRAINT fk_subscription_list_service_category FOREIGN KEY (service_category_uuid) REFERENCES service_category (service_category_uuid) ON DELETE CASCADE
    );

-- ======================
-- TABLE: transactions
-- ======================
CREATE TABLE
    transactions (
        transaction_uuid UUID PRIMARY KEY DEFAULT gen_random_uuid (),
        package_uuid UUID NOT NULL,
        owner_uuid UUID NOT NULL,
        service_category_uuid UUID NOT NULL,
        subscription_type VARCHAR(50) NOT NULL,
        transaction_date TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        amount DECIMAL(10, 2) NOT NULL,
        status VARCHAR(50) NOT NULL,
        duration INTEGER NOT NULL,
        duration_unit VARCHAR(20) NOT NULL, -- e.g., "day", "month", "year"
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        CONSTRAINT fk_transactions_package FOREIGN KEY (package_uuid) REFERENCES packages (package_uuid) ON DELETE CASCADE,
        CONSTRAINT fk_transactions_owner FOREIGN KEY (owner_uuid) REFERENCES owner (owner_uuid) ON DELETE CASCADE,
        CONSTRAINT fk_transactions_service_category FOREIGN KEY (service_category_uuid) REFERENCES service_category (service_category_uuid) ON DELETE CASCADE
    );

-- ======================
-- TABLE: subscription_payments
-- ======================
CREATE TABLE
    subscription_payment (
        subscription_payment_uuid UUID PRIMARY KEY DEFAULT gen_random_uuid (),
        subscription_payment_code VARCHAR(13) UNIQUE,
        transaction_uuid UUID NOT NULL,
        method VARCHAR(100) NOT NULL,
        status VARCHAR(50) NOT NULL,
        subscription_amount DECIMAL(10, 2) NOT NULL,
        tax_amount DECIMAL(10, 2) NULL,
        tax_percentage INT NULL,
        total_amount DECIMAL(10, 2) NOT NULL,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        CONSTRAINT fk_payment_transaction FOREIGN KEY (transaction_uuid) REFERENCES transactions (transaction_uuid) ON DELETE CASCADE
    );