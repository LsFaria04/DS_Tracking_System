--Remove any content that already exists in the db
DROP TABLE IF EXISTS orders CASCADE;
DROP TABLE IF EXISTS order_status_history CASCADE;
DROP table IF EXISTS orders_products CASCADE;
DROP TYPE IF EXISTS order_states CASCADE;

--Create an enum with all the possible order states
CREATE TYPE order_state AS ENUM ('PROCESSING', 'SHIPPED', 'IN TRANSIT', 'OUT FOR DELIVERY', 'CANCELLED', 'RETURNED', 'FAILED DELIVERY');

-- Orders table: stores basic order information
CREATE TABLE orders (
    id SERIAL PRIMARY KEY,
    customer_id INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    tracking_code TEXT UNIQUE NOT NULL,
    delivery_estimates date NOT NULL,
    delivery_address TEXT NOT NULL
);

-- Order status history: immutable audit log of status changes. It will be also stored in the blockchain
CREATE TABLE order_status_history (
    id SERIAL PRIMARY KEY,
    order_id INTEGER NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    order_status order_state NOT NULL ,
    timestamp_history TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    note TEXT,
    order_location TEXT NOT NULL
);

-- Contains the products that are part of the order
CREATE TABLE order_products (
    id SERIAL PRIMARY KEY,
    order_id INTEGER NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    product_id INTEGER NOT NULL,
    quantity INTEGER NOT NULL CHECK(quantity > 0)
);

-- Index for faster lookups
CREATE UNIQUE INDEX idx_orders_tracking_code ON orders(tracking_code);
CREATE INDEX idx_order_products_order_id ON order_products USING HASH(order_id);
CREATE INDEX idx_status_order_time ON order_status_history(order_id, timestamp_history); --when searching for a specific order id and sorting the data by 

--Triggers
CREATE OR REPLACE FUNCTION update_history_violation()
RETURNS TRIGGER AS $$
BEGIN
   RAISE EXCEPTION 'Updates and Deletes are not allowed on this table';
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_history
BEFORE UPDATE OR DELETE ON order_status_history
FOR EACH ROW
EXECUTE FUNCTION update_history_violation();
