-- Sample seed data for development and testing
-- This creates a few sample orders to demonstrate the tracking system

INSERT INTO orders (order_number, customer_id, current_status) VALUES
    ('ORD-2024-001', 'CUST-100', 'IN_PRODUCTION'),
    ('ORD-2024-002', 'CUST-101', 'SHIPPED'),
    ('ORD-2024-003', 'CUST-102', 'RECEIVED');

-- Initial status history entries
INSERT INTO order_status_history (order_id, status, notes) VALUES
    (1, 'IN_PRODUCTION', 'Order started production'),
    (2, 'IN_PRODUCTION', 'Order started production'),
    (2, 'READY_FOR_SHIPMENT', 'Production completed'),
    (2, 'SHIPPED', 'Package dispatched'),
    (3, 'IN_PRODUCTION', 'Order started production'),
    (3, 'READY_FOR_SHIPMENT', 'Production completed'),
    (3, 'SHIPPED', 'Package dispatched'),
    (3, 'RECEIVED', 'Delivered to customer');
