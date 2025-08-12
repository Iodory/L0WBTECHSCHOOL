CREATE TABLE orders (
    order_uid VARCHAR(50) PRIMARY KEY,
    track_number VARCHAR(100),
    entry VARCHAR(100),
    locale VARCHAR(20),
    internal_signature TEXT,
    customer_id VARCHAR(50),
    delivery_service VARCHAR(100),
    shardkey VARCHAR(50),
    sm_id INT,
    date_created TIMESTAMP,
    oof_shard VARCHAR(50)
);
