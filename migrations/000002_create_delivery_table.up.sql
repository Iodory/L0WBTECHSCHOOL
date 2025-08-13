CREATE TABLE IF NOT EXISTS delivery (
    order_uid VARCHAR(50) PRIMARY KEY,
    name VARCHAR(100),
    phone VARCHAR(20),
    zip VARCHAR(20),
    city VARCHAR(100),
    address TEXT,
    region VARCHAR(100),
    email VARCHAR(100),
    CONSTRAINT fk_delivery_order FOREIGN KEY(order_uid) REFERENCES orders(order_uid) ON DELETE CASCADE
);
