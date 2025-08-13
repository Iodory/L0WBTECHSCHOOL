CREATE TABLE IF NOT EXISTS payment (
    order_uid VARCHAR(50),
    transaction VARCHAR(100),
    request_id VARCHAR(100),
    currency VARCHAR(10),
    provider VARCHAR(100),
    amount INT,
    payment_dt BIGINT,
    bank VARCHAR(100),
    delivery_cost INT,
    goods_total INT,
    custom_fee INT,
    CONSTRAINT fk_payment_order FOREIGN KEY(order_uid) REFERENCES orders(order_uid)
);
