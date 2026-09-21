-- Data contoh untuk development.
-- Jalankan setelah migrations/000001_init_schema.up.sql.


CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- Password: admin -> admin123, sisanya -> password123
INSERT INTO users (name, username, password, role) VALUES
    ('Admin Lelang', 'admin', crypt('admin123',    gen_salt('bf', 10)), 'ADMIN'),
    ('John Doe',     'john',  crypt('password123', gen_salt('bf', 10)), 'USER'),
    ('Jane Smith',   'jane',  crypt('password123', gen_salt('bf', 10)), 'USER'),
    ('Mike Lee',     'mike',  crypt('password123', gen_salt('bf', 10)), 'USER'),
    ('Siti Rahma',   'siti',  crypt('password123', gen_salt('bf', 10)), 'USER');

INSERT INTO auctions (auction_name, description, image_link, bid_winner_id, starting_price, last_price, created_at, end_time, is_completed) VALUES
    ('iPhone 15 Pro',
     'Brand new iPhone 15 Pro 256GB Natural Titanium',
     'https://cdn.example.com/images/iphone15pro.jpg',
    NULL, 10000000, 16000000,
     NOW() - INTERVAL '2 days', NOW() + INTERVAL '7 days', FALSE),

    ('MacBook Air M3',
     'MacBook Air 13 inch M3 8GB/256GB Midnight, garansi resmi',
     'https://cdn.example.com/images/macbookair-m3.jpg',
    NULL, 15000000, 16250000,
     NOW() - INTERVAL '1 day', NOW() + INTERVAL '3 days', FALSE),

    ('Sepeda Brompton C Line',
     'Brompton C Line Explore 6-speed, kondisi mulus jarang dipakai',
     'https://cdn.example.com/images/brompton-cline.jpg',
    3, 25000000, 31000000,
     NOW() - INTERVAL '10 days', NOW() - INTERVAL '1 day', TRUE),

    ('Kamera Fujifilm X-T5',
     'Fujifilm X-T5 body only, shutter count rendah, fullset',
     'https://cdn.example.com/images/fujifilm-xt5.jpg',
    NULL, 18000000, 18000000,
     NOW() - INTERVAL '5 hours', NOW() + INTERVAL '14 days', FALSE),

    ('PlayStation 5 Slim',
     'PS5 Slim Digital Edition, segel belum dibuka',
     'https://cdn.example.com/images/ps5-slim.jpg',
    NULL, 7500000, 7800000,
     NOW() - INTERVAL '3 hours', NOW() + INTERVAL '5 days', FALSE),
     (
        'Apple Watch Series 8 45mm GPS', 
        'Midnight Aluminum Case, battery health 92%, mulus no dent', 
        'https://cdn.example.com/images/apple-watch-s8.jpg', 
        2, 
        4000000.00, 
        5100000.00, 
        NOW() - INTERVAL '14 days', 
        NOW() - INTERVAL '4 days', 
        TRUE
    ),
    (
        'Fujifilm X-T30 II Body Only', 
        'Warna Silver, SC rendah 3 ribuan, kelengkapan box fullset', 
        'https://cdn.example.com/images/fujifilm-xt30.jpg', 
        2, 
        10000000.00, 
        12300000.00, 
        NOW() - INTERVAL '20 days', 
        NOW() - INTERVAL '10 days', 
        TRUE
    ),
    (
        'iPad Air 5 M1 64GB Wi-Fi Space Gray', 
        'Layar bening terpasang paperlike screen protector, garansi aktif', 
        'https://cdn.example.com/images/ipad-air-5.jpg', 
        5, 
        7500000.00, 
        8900000.00, 
        NOW() - INTERVAL '10 days', 
        NOW() - INTERVAL '2 days', 
        TRUE
    ),
    (
        'LG DualUp Monitor 28SD750 16:18 Ergo Stand', 
        'Layar unik SDQHD IPS, kondisi normal no dead pixel, lengkap box', 
        'https://cdn.example.com/images/lg-dualup.jpg', 
        4, 
        6500000.00, 
        7600000.00, 
        NOW() - INTERVAL '18 days', 
        NOW() - INTERVAL '7 days', 
        TRUE
    ),
    (
        'Samsung Galaxy S23 Ultra 12/256GB', 
        'Warna Phantom Black, garansi Sein resmi, pemakaian terawat', 
        'https://cdn.example.com/images/s23-ultra.jpg', 
        3, 
        11000000.00, 
        13200000.00, 
        NOW() - INTERVAL '9 days', 
        NOW() - INTERVAL '1 day', 
        TRUE
    ),
    (
        'Nintendo Switch OLED Zelda Tears of the Kingdom Edition', 
        'Fullset mulus, terpasang tempered glass, joycon no drift', 
        'https://cdn.example.com/images/switch-oled-totk.jpg', 
        2, 
        4200000.00, 
        5000000.00, 
        NOW() - INTERVAL '25 days', 
        NOW() - INTERVAL '15 days', 
        TRUE
    ),
    (
        'Bose QuietComfort 45 Headphone', 
        'Warna Smoke White, earpad masih empuk, fungsi normal 100%', 
        'https://cdn.example.com/images/bose-qc45.jpg', 
        5, 
        2800000.00, 
        3400000.00, 
        NOW() - INTERVAL '16 days', 
        NOW() - INTERVAL '6 days', 
        TRUE
    ),
    (
        'DJI Mini 3 Pro Fly More Combo', 
        'Drone mulus no crash, sensor lancar, dapat 3 baterai', 
        'https://cdn.example.com/images/dji-mini-3.jpg', 
        4, 
        9500000.00, 
        11800000.00, 
        NOW() - INTERVAL '30 days', 
        NOW() - INTERVAL '18 days', 
        TRUE
    ),
    (
        'Dyson V12 Detect Slim Vacuum', 
        'Kondisi fisik 90%, hisapan kuat, aksesoris lengkap', 
        'https://cdn.example.com/images/dyson-v12.jpg', 
        3, 
        6000000.00, 
        7200000.00, 
        NOW() - INTERVAL '11 days', 
        NOW() - INTERVAL '3 days', 
        TRUE
    ),
    (
        'ASUS ROG Ally Z1 Extreme 512GB', 
        'Handheld gaming PC, terpasang SSD 1TB upgrade, garansi resmi', 
        'https://cdn.example.com/images/rog-ally.jpg', 
        2, 
        7000000.00, 
        8500000.00, 
        NOW() - INTERVAL '7 days', 
        NOW() - INTERVAL '1 day', 
        TRUE
    ),
    (
        'iPhone 15 Pro Max 256GB Natural Titanium', 
        'Garansi resmi iBox aktif, kondisi 99% mulus like new, battery health 100%', 
        'https://cdn.example.com/images/iphone-15-promax.jpg', 
        NULL, 16000000.00, 
        17500000.00, 
        NOW() - INTERVAL '2 days', 
        NOW() + INTERVAL '3 days', 
        FALSE
    ),
    (
        'Sony PlayStation Portal Remote Player', 
        'Kondisi gress pemakaian seminggu, terpasang tempered glass, lengkap box', 
        'https://cdn.example.com/images/ps-portal.jpg', 
        NULL, 3200000.00, 
        3600000.00, 
        NOW() - INTERVAL '1 day', 
        NOW() + INTERVAL '5 days', 
        FALSE
    ),
    (
        'MacMini M2 Pro 16/512GB', 
        'Pembelian 2024, performa mantap untuk video editing 4K, fullset original', 
        'https://cdn.example.com/images/macmini-m2.jpg', 
        NULL, 15000000.00, 
        15800000.00, 
        NOW() - INTERVAL '3 days', 
        NOW() + INTERVAL '1 day', 
        FALSE
    ),
    (
        'Garmin Forerunner 965 Black', 
        'Smartwatch lari AMOLED, bezel titanium, battery awet 10 hari, sensor akurat', 
        'https://cdn.example.com/images/garmin-fr965.jpg', 
        NULL, 7500000.00, 
        8200000.00, 
        NOW() - INTERVAL '12 hours', 
        NOW() + INTERVAL '2 days', 
        FALSE
    ),
    (
        'Logitech MX Master 3S Wireless Mouse', 
        'Warna Pale Gray, silent click, kondisi mulus no minus, garansi distributor', 
        'https://cdn.example.com/images/mx-master-3s.jpg', 
        NULL, 1100000.00, 
        1350000.00, 
        NOW() - INTERVAL '6 hours', 
        NOW() + INTERVAL '4 days', 
        FALSE
    );

INSERT INTO bids (auction_id, user_id, bid_price, created_at)
SELECT a.id, u.id, v.bid_price, NOW() - v.age
FROM (VALUES
    ('iPhone 15 Pro',          'mike', 14000000::NUMERIC, INTERVAL '30 hours'),
    ('iPhone 15 Pro',          'jane', 14500000::NUMERIC, INTERVAL '26 hours'),
    ('iPhone 15 Pro',          'john', 15000000::NUMERIC, INTERVAL '20 hours'),
    ('iPhone 15 Pro',          'siti', 16000000::NUMERIC, INTERVAL '4 hours'),

    ('MacBook Air M3',         'john', 15500000::NUMERIC, INTERVAL '18 hours'),
    ('MacBook Air M3',         'siti', 16250000::NUMERIC, INTERVAL '9 hours'),

    ('Sepeda Brompton C Line', 'jane', 26000000::NUMERIC, INTERVAL '8 days'),
    ('Sepeda Brompton C Line', 'mike', 27500000::NUMERIC, INTERVAL '6 days'),
    ('Sepeda Brompton C Line', 'john', 30000000::NUMERIC, INTERVAL '4 days'),
    ('Sepeda Brompton C Line', 'jane', 31000000::NUMERIC, INTERVAL '2 days'),

    ('PlayStation 5 Slim',     'mike',  7800000::NUMERIC, INTERVAL '2 hours'),

    ('Apple Watch Series 8 45mm GPS', 'john', 4700000::NUMERIC, INTERVAL '12 days'),
    ('Apple Watch Series 8 45mm GPS', 'john', 5100000::NUMERIC, INTERVAL '5 days'),

    ('Fujifilm X-T30 II Body Only', 'mike', 11000000::NUMERIC, INTERVAL '18 days'),
    ('Fujifilm X-T30 II Body Only', 'john', 12300000::NUMERIC, INTERVAL '11 days'),

    ('iPad Air 5 M1 64GB Wi-Fi Space Gray', 'jane', 8000000::NUMERIC, INTERVAL '8 days'),
    ('iPad Air 5 M1 64GB Wi-Fi Space Gray', 'siti', 8900000::NUMERIC, INTERVAL '3 days'),

    ('LG DualUp Monitor 28SD750 16:18 Ergo Stand', 'john', 7000000::NUMERIC, INTERVAL '15 days'),
    ('LG DualUp Monitor 28SD750 16:18 Ergo Stand', 'mike', 7600000::NUMERIC, INTERVAL '8 days'),

    ('Samsung Galaxy S23 Ultra 12/256GB', 'siti', 12000000::NUMERIC, INTERVAL '7 days'),
    ('Samsung Galaxy S23 Ultra 12/256GB', 'jane', 13200000::NUMERIC, INTERVAL '2 days'),

    ('Nintendo Switch OLED Zelda Tears of the Kingdom Edition', 'mike', 4500000::NUMERIC, INTERVAL '22 days'),
    ('Nintendo Switch OLED Zelda Tears of the Kingdom Edition', 'john', 5000000::NUMERIC, INTERVAL '16 days'),

    ('Bose QuietComfort 45 Headphone', 'jane', 3000000::NUMERIC, INTERVAL '14 days'),
    ('Bose QuietComfort 45 Headphone', 'siti', 3400000::NUMERIC, INTERVAL '7 days'),

    ('DJI Mini 3 Pro Fly More Combo', 'john', 10500000::NUMERIC, INTERVAL '25 days'),
    ('DJI Mini 3 Pro Fly More Combo', 'mike', 11800000::NUMERIC, INTERVAL '19 days'),

    ('Dyson V12 Detect Slim Vacuum', 'siti', 6500000::NUMERIC, INTERVAL '9 days'),
    ('Dyson V12 Detect Slim Vacuum', 'jane', 7200000::NUMERIC, INTERVAL '4 days'),

    ('ASUS ROG Ally Z1 Extreme 512GB', 'mike', 7800000::NUMERIC, INTERVAL '6 days'),
    ('ASUS ROG Ally Z1 Extreme 512GB', 'john', 8500000::NUMERIC, INTERVAL '2 days'),

    ('iPhone 15 Pro Max 256GB Natural Titanium', 'siti', 16500000::NUMERIC, INTERVAL '1 day'),
    ('iPhone 15 Pro Max 256GB Natural Titanium', 'john', 17500000::NUMERIC, INTERVAL '4 hours'),

    ('Sony PlayStation Portal Remote Player', 'mike', 3400000::NUMERIC, INTERVAL '20 hours'),
    ('Sony PlayStation Portal Remote Player', 'jane', 3600000::NUMERIC, INTERVAL '6 hours'),

    ('MacMini M2 Pro 16/512GB', 'john', 15300000::NUMERIC, INTERVAL '2 days'),
    ('MacMini M2 Pro 16/512GB', 'siti', 15800000::NUMERIC, INTERVAL '8 hours'),

    ('Garmin Forerunner 965 Black', 'jane', 7800000::NUMERIC, INTERVAL '10 hours'),
    ('Garmin Forerunner 965 Black', 'mike', 8200000::NUMERIC, INTERVAL '3 hours'),

    ('Logitech MX Master 3S Wireless Mouse', 'siti', 1200000::NUMERIC, INTERVAL '5 hours'),
    ('Logitech MX Master 3S Wireless Mouse', 'john', 1350000::NUMERIC, INTERVAL '1 hour')
) AS v(auction_name, username, bid_price, age)
JOIN auctions a ON a.auction_name = v.auction_name
JOIN users u ON u.username = v.username;