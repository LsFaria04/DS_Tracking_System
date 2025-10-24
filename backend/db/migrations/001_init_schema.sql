-- Initial schema for order tracking system
-- This migration creates the foundational tables for tracking order status

-- Orders table: stores basic order information
CREATE TABLE IF NOT EXISTS orders (
    id SERIAL PRIMARY KEY,
    order_number VARCHAR(50) UNIQUE NOT NULL,
    customer_id VARCHAR(100) NOT NULL,
    current_status VARCHAR(50) NOT NULL DEFAULT 'IN_PRODUCTION',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Order status history: immutable audit log of status changes
CREATE TABLE IF NOT EXISTS order_status_history (
    id SERIAL PRIMARY KEY,
    order_id INTEGER NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    status VARCHAR(50) NOT NULL,
    timestamp TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    notes TEXT
);

-- Index for faster lookups
CREATE INDEX idx_orders_order_number ON orders(order_number);
CREATE INDEX idx_orders_customer_id ON orders(customer_id);
CREATE INDEX idx_status_history_order_id ON order_status_history(order_id);

-- Valid status values (enforced at application level for now):
-- 'IN_PRODUCTION', 'READY_FOR_SHIPMENT', 'SHIPPED', 'RECEIVED'
