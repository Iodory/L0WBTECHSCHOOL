CREATE TABLE items (
    id SERIAL PRIMARY KEY,       
    order_uid VARCHAR(50),
    chrt_id BIGINT,
    track_number VARCHAR(100),
    price INT,
    rid VARCHAR(100),
    name VARCHAR(255),
    sale INT,
    size VARCHAR(50),
    total_price INT,
    nm_id BIGINT,
    brand VARCHAR(100),
    status INT,
    CONSTRAINT fk_items_order FOREIGN KEY(order_uid) REFERENCES orders(order_uid) ON DELETE CASCADE
);
