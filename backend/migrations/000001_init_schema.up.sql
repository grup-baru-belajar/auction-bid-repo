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

