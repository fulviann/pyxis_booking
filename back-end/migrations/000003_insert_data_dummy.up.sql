-- ======================
-- INSERT: merchants
-- ======================
INSERT INTO merchants (merchant_code, name, address, phone_number, email)
VALUES
('MRC15042026001', 'GOR Basket Jaya', 'Jl. Raya Surabaya No. 1', '081234567890', 'gorjaya@gmail.com'),
('MRC15042026002', 'Arena Basket Pro', 'Jl. Ahmad Yani No. 45', '081298765432', 'arena.pro@gmail.com'),
('MRC15042026003', 'Lapangan Serbaguna Sport Center', 'Jl. Diponegoro No. 10', '082112223334', 'sportcenter@gmail.com');

-- ======================
-- INSERT: court (BASKET)
-- ======================
INSERT INTO court (court_code, merchant_uuid, name, type, price_per_hour, status, capacity)
VALUES
('CRT15042026001',
 (SELECT merchant_uuid FROM merchants WHERE merchant_code = 'MRC15042026001'),
 'Lapangan Basket A', 'Full Court Indoor', 200000, 'AVAILABLE', 10),

('CRT15042026002',
 (SELECT merchant_uuid FROM merchants WHERE merchant_code = 'MRC15042026001'),
 'Lapangan Basket B', 'Half Court Outdor', 150000, 'AVAILABLE', 10),

('CRT15042026003',
 (SELECT merchant_uuid FROM merchants WHERE merchant_code = 'MRC15042026001'),
 'Lapangan Basket Premium', 'Full Court Indoor', 250000, 'AVAILABLE', 10),

 ('CRT15042026004',
 (SELECT merchant_uuid FROM merchants WHERE merchant_code = 'MRC15042026002'),
 'Lapangan Basket A', 'Full Court Indoor', 200000, 'AVAILABLE', 10),

('CRT15042026005',
 (SELECT merchant_uuid FROM merchants WHERE merchant_code = 'MRC15042026002'),
 'Lapangan Basket B', 'Half Court Outdor', 150000, 'AVAILABLE', 10),

('CRT15042026006',
 (SELECT merchant_uuid FROM merchants WHERE merchant_code = 'MRC15042026002'),
 'Lapangan Basket Premium', 'Full Court Indoor', 250000, 'AVAILABLE', 10),
 
 ('CRT15042026007',
 (SELECT merchant_uuid FROM merchants WHERE merchant_code = 'MRC15042026003'),
 'Lapangan Basket A', 'Full Court Indoor', 200000, 'AVAILABLE', 10),

('CRT15042026008',
 (SELECT merchant_uuid FROM merchants WHERE merchant_code = 'MRC15042026003'),
 'Lapangan Basket B', 'Half Court Outdor', 150000, 'AVAILABLE', 10),

('CRT15042026009',
 (SELECT merchant_uuid FROM merchants WHERE merchant_code = 'MRC15042026003'),
 'Lapangan Basket Premium', 'Full Court Indoor', 250000, 'AVAILABLE', 10);

-- ======================
-- INSERT: member password hash: '123456' -> $2a$10$7EqJtq98hPqEX7fNZaFWoOHiY3VYj6lZ9E9QnUnxE27XGr0Ccs9xG
-- ======================
 INSERT INTO member (
    member_code, name, phone_number, email, password, position
)
VALUES
('MBR15042026001', 'Susi Susanti', '081234567890', 'susi@gmail.com',
 '$2a$10$7EqJtq98hPqEX7fNZaFWoOHiY3VYj6lZ9E9QnUnxE27XGr0Ccs9xG', 'USER'),

('MBR15042026002', 'Budi Santoso', '081298765432', 'budi@gmail.com',
 '$2a$10$7EqJtq98hPqEX7fNZaFWoOHiY3VYj6lZ9E9QnUnxE27XGr0Ccs9xG', 'USER'),

('MBR15042026003', 'Admin System', '080000000000', 'admin@system.com',
 '$2a$10$7EqJtq98hPqEX7fNZaFWoOHiY3VYj6lZ9E9QnUnxE27XGr0Ccs9xG', 'ADMIN');

-- ======================
-- INSERT: blackouts
-- ======================
INSERT INTO blackouts (
    court_uuid,
    merchant_uuid,
    title,
    start_date,
    end_date,
    is_full_day
)
VALUES
-- 1. Maintenance lapangan tertentu
(
    (SELECT court_uuid FROM court WHERE court_code = 'CRT15042026001'),
    NULL,
    'Maintenance Lapangan',
    '2026-04-20 08:00:00+07',
    '2026-04-20 12:00:00+07',
    FALSE
),

-- 2. Libur nasional (semua lapangan di merchant)
(
    NULL,
    (SELECT merchant_uuid FROM merchants WHERE merchant_code = 'MRC15042026001'),
    'Libur Hari Raya',
    '2026-04-21 00:00:00+07',
    '2026-04-21 23:59:59+07',
    TRUE
),

-- 3. Event khusus (blok 1 lapangan full day)
(
    (SELECT court_uuid FROM court WHERE court_code = 'CRT15042026003'),
    NULL,
    'Turnamen Basket Lokal',
    '2026-04-22 00:00:00+07',
    '2026-04-22 23:59:59+07',
    TRUE
);