CREATE TABLE orders
(
    id                SERIAL PRIMARY KEY,
    external_order_id VARCHAR,
    status            VARCHAR,
    price             VARCHAR,
    quantity          VARCHAR,
    side              VARCHAR,
    raw_request       VARCHAR,
    raw_answer        VARCHAR,
    created_at        TIMESTAMP DEFAULT CURRENT_TIMESTAMP
)