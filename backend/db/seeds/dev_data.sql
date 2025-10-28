-- Insert sample products
INSERT INTO products (id, name, price)
VALUES
  (201, 'Gaming Headset HYPERX Cloud II', 62.99),
  (202, 'Mechanical Keyboard Logitech G Pro', 67.98),
  (203, 'Gaming Mouse RAZER Deathadder V2', 49.99),
  (204, 'Mousepad LOGITECH G240', 17.99),
  (205, 'USB-C Cable 2m', 5.99);

-- Insert sample orders
INSERT INTO orders (customer_id, tracking_code, delivery_estimate, delivery_address)
VALUES
  (101, 'TRACK001', CURRENT_DATE + INTERVAL '5 days', '123 Elm Street'),
  (102, 'TRACK002', CURRENT_DATE + INTERVAL '3 days', '456 Oak Avenue'),
  (103, 'TRACK003', CURRENT_DATE + INTERVAL '7 days', '789 Pine Road');

-- Insert status history for each order
INSERT INTO order_status_history (order_id, order_status, note, order_location)
VALUES
  (1, 'PROCESSING', 'Order received', 'Warehouse A'),
  (1, 'SHIPPED', 'Left warehouse', 'Warehouse A'),
  (2, 'PROCESSING', 'Order received', 'Warehouse B'),
  (2, 'IN TRANSIT', 'On the way', 'Distribution Center'),
  (3, 'PROCESSING', 'Order received', 'Warehouse C'),
  (3, 'CANCELLED', 'Customer requested cancellation', 'Warehouse C');

-- Insert products for each order
INSERT INTO order_products (order_id, product_id, quantity)
VALUES
  (1, 201, 2),
  (1, 202, 1),
  (2, 203, 3),
  (3, 204, 1),
  (3, 205, 2);

