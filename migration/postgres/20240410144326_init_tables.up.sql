CREATE TABLE IF NOT EXISTS pools (
    id SERIAL PRIMARY KEY,
    address VARCHAR(64) NOT NULL,
    start_block BIGINT NOT NULL,
    token0 VARCHAR(64) NOT NULL,
    token0_decimals INT NOT NULL,
    token1 VARCHAR(64) NOT NULL,
    token1_decimals INT NOT NULL,
    created_at BIGINT NOT NULL,
    UNIQUE (address)
);

CREATE TABLE IF NOT EXISTS txs (
    id SERIAL PRIMARY KEY,
    pool_id INT NOT NULL,
    tx_hash VARCHAR(128) NOT NULL,
    block_number BIGINT NOT NULL,
    block_time BIGINT NOT NULL,
    gas BIGINT NOT NULL,
    gas_price DECIMAL(80, 0) NOT NULL,
    receipt JSON NOT NULL,
    is_finalized BOOLEAN NOT NULL,
    created_at BIGINT NOT NULL,
    FOREIGN KEY (pool_id) REFERENCES pools(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS swap_events (
    id SERIAL PRIMARY KEY,
    pool_id INT NOT NULL,
    tx_id INT NOT NULL,
    amount0 VARCHAR(128),
    amount1 VARCHAR(128),
    price VARCHAR(128),
    created_at BIGINT NOT NULL,
    FOREIGN KEY (pool_id) REFERENCES pools(id) ON DELETE CASCADE,
    FOREIGN KEY (tx_id) REFERENCES txs(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS ethusdt_klines (
    id SERIAL PRIMARY KEY,
    open_time BIGINT NOT NULL,
    close_time BIGINT NOT NULL,
    open_price VARCHAR(32),
    high_price VARCHAR(32),
    low_price VARCHAR(32),
    close_price VARCHAR(32),
    ohlc4 DECIMAL(20, 10),
    created_at BIGINT NOT NULL,
    UNIQUE (open_time)
);

CREATE TABLE IF NOT EXISTS block_cursors (
    id SERIAL PRIMARY KEY,
    pool_id INT NOT NULL,
    type VARCHAR(16),
    block_number BIGINT NOT NULL,
    extra JSON,
    created_at BIGINT NOT NULL,
    updated_at BIGINT NOT NULL,
    FOREIGN KEY (pool_id) REFERENCES pools(id) ON DELETE CASCADE,
    UNIQUE (pool_id, type)
);