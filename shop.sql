CREATE TABLE orders (
                        order_id SERIAL PRIMARY KEY,
                        status VARCHAR(50)
);

CREATE TABLE inventory (
                           product_id SERIAL PRIMARY KEY,
                           product_name VARCHAR(100),
                           quantity INT
);

CREATE TABLE payments (
                          payment_id SERIAL PRIMARY KEY,
                          order_id INT REFERENCES orders(order_id),
                          amount DECIMAL(10, 2)
);

CREATE TABLE deliveries (
                            delivery_id SERIAL PRIMARY KEY,
                            order_id INT REFERENCES orders(order_id),
                            address VARCHAR(255),
                            status VARCHAR(50)
);

INSERT INTO inventory (product_name, quantity) VALUES
                                                   ('Product A', 100),
                                                   ('Product B', 200),
                                                   ('Product C', 150);
