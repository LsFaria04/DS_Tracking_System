-- Insert storage locations with real GPS coordinates
INSERT INTO storages (name, address, latitude, longitude) VALUES
('Main Warehouse Lisboa', 'Av. da Liberdade, 1250-096 Lisboa, Portugal', 38.7223, -9.1393),
('Distribution Center Porto', 'Rua de Santa Catarina, 4000-442 Porto, Portugal', 41.1496, -8.6109),
('Regional Hub Coimbra', 'Praça da República, 3000-343 Coimbra, Portugal', 40.2033, -8.4103),
('Logistics Center Faro', 'Av. da República, 8000-078 Faro, Portugal', 37.0194, -7.9304),
('Distribution Hub Braga', 'Praça da República, 4710-305 Braga, Portugal', 41.5454, -8.4265),
('Regional Center Évora', 'Praça do Giraldo, 7000-508 Évora, Portugal', 38.5714, -7.9087),
('Logistics Hub Aveiro', 'Av. Dr. Lourenço Peixinho, 3800-167 Aveiro, Portugal', 40.6443, -8.6455),
('Distribution Center Setúbal', 'Av. Luísa Todi, 2900-456 Setúbal, Portugal', 38.5244, -8.8882),
('Regional Hub Viseu', 'Rua Formosa, 3500-161 Viseu, Portugal', 40.6566, -7.9122),
('Madeira Hub Funchal', 'Av. Arriaga, 9000-060 Funchal, Madeira', 32.6669, -16.9241),
('Açores Hub Ponta Delgada', 'Av. Infante Dom Henrique, 9500-150 Ponta Delgada, Açores', 37.7412, -25.6756);

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
  (101, 'TRACK001', CURRENT_DATE + INTERVAL '5 days', 'Rua Augusta 123, 1100-053 Lisboa, Portugal'),
  (102, 'TRACK002', CURRENT_DATE + INTERVAL '3 days', 'Rua das Flores 456, 4050-265 Porto, Portugal'),
  (103, 'TRACK003', CURRENT_DATE + INTERVAL '7 days', 'Rua Ferreira Borges 789, 3000-180 Coimbra, Portugal'),
  (104, 'TRACK004', CURRENT_DATE + INTERVAL '10 days', 'Rua Dr. Fernão Ornelas 321, 9000-082 Funchal, Madeira, Portugal'),
  (105, 'TRACK005', CURRENT_DATE + INTERVAL '12 days', 'Rua Machado dos Santos 654, 9500-092 Ponta Delgada, Açores, Portugal');

-- Insert status history for each order
INSERT INTO order_status_history (order_id, order_status, note, order_location, storage_id)
VALUES
  -- Order 1: Lisboa route
  (1, 'PROCESSING', 'Order received', 'Main Warehouse Lisboa', 1),
  (1, 'SHIPPED', 'Left warehouse', 'Main Warehouse Lisboa', 1),
  
  -- Order 2: Porto route
  (2, 'PROCESSING', 'Order received', 'Distribution Center Porto', 2),
  (2, 'IN TRANSIT', 'On the way', 'Regional Hub Coimbra', 3),
  
  -- Order 3: Cancelled
  (3, 'PROCESSING', 'Order received', 'Logistics Center Faro', 4),
  (3, 'CANCELLED', 'Customer requested cancellation', 'Logistics Center Faro', 4),
  
  -- Order 4: Madeira route (Lisboa -> Madeira)
  (4, 'PROCESSING', 'Order received', 'Main Warehouse Lisboa', 1),
  (4, 'SHIPPED', 'Left warehouse', 'Main Warehouse Lisboa', 1),
  (4, 'IN TRANSIT', 'In transit to Madeira', 'Main Warehouse Lisboa', 1),
  (4, 'IN TRANSIT', 'Arrived at Madeira hub', 'Madeira Hub Funchal', 10),
  (4, 'OUT FOR DELIVERY', 'Out for delivery in Funchal', 'Madeira Hub Funchal', 10),
  
  -- Order 5: Açores route (Porto -> Açores)
  (5, 'PROCESSING', 'Order received', 'Distribution Center Porto', 2),
  (5, 'SHIPPED', 'Left warehouse', 'Distribution Center Porto', 2),
  (5, 'IN TRANSIT', 'In transit to Açores', 'Distribution Center Porto', 2),
  (5, 'IN TRANSIT', 'Arrived at Açores hub', 'Açores Hub Ponta Delgada', 11);

-- Insert products for each order
INSERT INTO order_products (order_id, product_id, quantity)
VALUES
  (1, 201, 2),
  (1, 202, 1),
  (2, 203, 3),
  (3, 204, 1),
  (3, 205, 2);

