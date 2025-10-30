--Remove any content that already exists in the db
DROP TABLE IF EXISTS orders CASCADE;
DROP TABLE IF EXISTS order_status_history CASCADE;
DROP TABLE IF EXISTS products CASCADE;
DROP TABLE IF EXISTS order_products CASCADE;
DROP TABLE IF EXISTS storages CASCADE;
DROP TYPE IF EXISTS order_state CASCADE;

--Create an enum with all the possible order states
CREATE TYPE order_state AS ENUM ('PROCESSING', 'SHIPPED', 'IN TRANSIT', 'OUT FOR DELIVERY', 'CANCELLED', 'RETURNED', 'FAILED DELIVERY');

-- Storages table: contains warehouse/storage locations with GPS coordinates
CREATE TABLE storages (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    address TEXT NOT NULL,
    latitude DECIMAL(10, 8) NOT NULL,
    longitude DECIMAL(11, 8) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- order table: stores basic order information
CREATE TABLE orders (
    id SERIAL PRIMARY KEY,
    customer_id INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    tracking_code TEXT UNIQUE NOT NULL,
    delivery_estimate date NOT NULL,
    delivery_address TEXT NOT NULL,
    price DECIMAL(10,2) NOT NULL DEFAULT 0.00 -- total price of the order
);

-- Order status history: immutable audit log of status changes. It will be also stored in the blockchain
CREATE TABLE order_status_history (
    id SERIAL PRIMARY KEY,
    order_id INTEGER NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    order_status order_state NOT NULL ,
    timestamp_history TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    note TEXT,
    order_location TEXT NOT NULL,
    storage_id INTEGER REFERENCES storages(id) ON DELETE SET NULL
);

-- Products table: contains available products for order
CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    price DECIMAL(10,2) NOT NULL
);

-- Contains the products that are part of the order
CREATE TABLE order_products (
    id SERIAL PRIMARY KEY,
    order_id INTEGER NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    product_id INTEGER NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    quantity INTEGER NOT NULL CHECK(quantity > 0)
);

-- Index for faster lookups
CREATE UNIQUE INDEX idx_order_tracking_code ON orders(tracking_code);
CREATE INDEX idx_order_product_order_id ON order_products USING HASH(order_id);
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

-- Trigger to update total price in order when order_products change
CREATE OR REPLACE FUNCTION update_order_price()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE orders
    SET price = (
        SELECT COALESCE(SUM(p.price * op.quantity),0)
        FROM order_products op
        JOIN products p ON op.product_id = p.id
        WHERE op.order_id = NEW.order_id
    )
    WHERE id = NEW.order_id;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_update_order_price
AFTER INSERT OR UPDATE OR DELETE ON order_products
FOR EACH ROW
EXECUTE FUNCTION update_order_price();
