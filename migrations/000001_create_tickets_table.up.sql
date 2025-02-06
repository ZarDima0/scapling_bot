CREATE TABLE tickets
(
    id            SERIAL PRIMARY KEY,
    symbol        VARCHAR,
    last_price    VARCHAR(255),
    high_price_24 VARCHAR(255) NOT NULL,
    low_price_24  VARCHAR(255) NOT NULL,
    volume_24     VARCHAR(255) NOT NULL,
    turnover_24   VARCHAR(255) NOT NULL,
    raw_ticket    JSON,
    created_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);