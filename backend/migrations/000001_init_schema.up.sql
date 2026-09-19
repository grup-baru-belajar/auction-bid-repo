CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    username VARCHAR(100) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    role VARCHAR(20) NOT NULL CHECK (role IN ('ADMIN', 'USER'))
);

CREATE TABLE auctions (
    id BIGSERIAL PRIMARY KEY,
    auction_name VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    image_link TEXT NOT NULL,
    bid_winner_id BIGINT NULL,
    starting_price NUMERIC(18,2) NOT NULL CHECK (starting_price > 0),
    last_price NUMERIC(18,2) NOT NULL CHECK (last_price >= starting_price),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    end_time TIMESTAMPTZ NOT NULL,
    is_completed BOOLEAN NOT NULL DEFAULT FALSE,

    CONSTRAINT fk_auction_bid_winner
        FOREIGN KEY (bid_winner_id)
        REFERENCES users(id)
        ON DELETE SET NULL
);

CREATE TABLE bids (
    id BIGSERIAL PRIMARY KEY,
    auction_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    bid_price NUMERIC(18,2) NOT NULL CHECK (bid_price > 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_bid_auction
        FOREIGN KEY (auction_id)
        REFERENCES auctions(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_bid_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE
);

CREATE INDEX idx_auctions_completed_end_time ON auctions(is_completed, end_time);

CREATE INDEX idx_bids_auction_price_desc ON bids(auction_id, bid_price DESC);


-- ============================================================
-- DUMMY DATA (development only) — separate later
-- ============================================================

CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- Password: admin -> admin123, sisanya -> password123
INSERT INTO users (name, username, password, role) VALUES
    ('Admin Lelang', 'admin', crypt('admin123',    gen_salt('bf', 10)), 'ADMIN'),
    ('John Doe',     'john',  crypt('password123', gen_salt('bf', 10)), 'USER'),
    ('Jane Smith',   'jane',  crypt('password123', gen_salt('bf', 10)), 'USER'),
    ('Mike Lee',     'mike',  crypt('password123', gen_salt('bf', 10)), 'USER'),
    ('Siti Rahma',   'siti',  crypt('password123', gen_salt('bf', 10)), 'USER');

INSERT INTO auctions (auction_name, description, image_link, starting_price, last_price, created_at, end_time, is_completed) VALUES
    ('iPhone 15 Pro',
     'Brand new iPhone 15 Pro 256GB Natural Titanium',
     'https://cdn.example.com/images/iphone15pro.jpg',
     10000000, 10000000,
     NOW() - INTERVAL '2 days', NOW() + INTERVAL '7 days', FALSE),

    ('MacBook Air M3',
     'MacBook Air 13 inch M3 8GB/256GB Midnight, garansi resmi',
     'https://cdn.example.com/images/macbookair-m3.jpg',
     15000000, 15000000,
     NOW() - INTERVAL '1 day', NOW() + INTERVAL '3 days', FALSE),

    ('Sepeda Brompton C Line',
     'Brompton C Line Explore 6-speed, kondisi mulus jarang dipakai',
     'https://cdn.example.com/images/brompton-cline.jpg',
     25000000, 25000000,
     NOW() - INTERVAL '10 days', NOW() - INTERVAL '1 day', TRUE),

    ('Kamera Fujifilm X-T5',
     'Fujifilm X-T5 body only, shutter count rendah, fullset',
     'https://cdn.example.com/images/fujifilm-xt5.jpg',
     18000000, 18000000,
     NOW() - INTERVAL '5 hours', NOW() + INTERVAL '14 days', FALSE),

    ('PlayStation 5 Slim',
     'PS5 Slim Digital Edition, segel belum dibuka',
     'https://cdn.example.com/images/ps5-slim.jpg',
     7500000, 7500000,
     NOW() - INTERVAL '3 hours', NOW() + INTERVAL '5 days', FALSE);

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

    ('PlayStation 5 Slim',     'mike',  7800000::NUMERIC, INTERVAL '2 hours')
) AS v(auction_name, username, bid_price, age)
JOIN auctions a ON a.auction_name = v.auction_name
JOIN users u ON u.username = v.username;

-- Sinkronkan last_price & bid_winner_id dengan bid tertinggi
UPDATE auctions a
SET last_price    = t.bid_price,
    bid_winner_id = t.user_id
FROM (
    SELECT DISTINCT ON (auction_id) auction_id, user_id, bid_price
    FROM bids
    ORDER BY auction_id, bid_price DESC, created_at ASC
) t
WHERE a.id = t.auction_id;