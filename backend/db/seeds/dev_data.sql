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
('Açores Hub Ponta Delgada', 'Av. Infante Dom Henrique, 9500-150 Ponta Delgada, Açores', 37.7412, -25.6756),
('Lisbon Airport (Cargo)', 'Aeroporto Humberto Delgado, 1700-111 Lisboa, Portugal', 38.7742, -9.1342),
('Porto Airport (Cargo)', 'Aeroporto Francisco Sá Carneiro, 4470-558 Maia, Portugal', 41.2481, -8.6814),
('Funchal Airport (Cargo)', 'Aeroporto da Madeira, 9100-105 Santa Cruz, Madeira', 32.6979, -16.7745),
('Ponta Delgada Airport (Cargo)', 'Aeroporto João Paulo II, 9500-749 Ponta Delgada, Açores', 37.7412, -25.6980),
('Faro Airport (Cargo)', 'Aeroporto Gago Coutinho, 8006-901 Faro, Portugal', 37.0142, -7.9659),
('Beja Airport (Cargo)', 'Aeroporto de Beja, EM528-2, 7800-745 Beja, Portugal', 38.0775, -7.9317);

-- Insert sample products
INSERT INTO products (id, name, price)
VALUES
  (201, 'Handcrafted Cork Wallet', 24.99),
  (202, 'Portuguese Ceramic Tile Set (6 pieces)', 45.50),
  (203, 'Artisan Olive Wood Cutting Board', 38.75),
  (204, 'Hand-painted Azulejo Coaster Set', 19.99),
  (205, 'Traditional Filigree Silver Earrings', 67.00);

-- Insert sample orders
INSERT INTO orders (customer_id, seller_id, seller_address, seller_latitude, seller_longitude, tracking_code, delivery_estimate, delivery_address, delivery_latitude, delivery_longitude)
VALUES
  -- Order 1: Seller near Lisboa (Almada) -> Main Warehouse Lisboa -> Delivery in Lisboa (Parque das Nações)
  (101, 501, 'Dona Lurdes, Rua Dom Afonso Henriques 12, 2800-012 Almada, Portugal', 38.6780, -9.1580, 'TRACK001', CURRENT_DATE + INTERVAL '5 days', 'Rua Padre Joaquim Alves Correia 5, 1990-152 Lisboa, Portugal', 38.7680, -9.1000),
  
  -- Order 2: Seller near Porto (Matosinhos) -> Distribution Center Porto -> Delivery in Coimbra (Santa Clara)
  (102, 502, 'Dona Maria, Rua Brito Capelo 100, 4450-073 Matosinhos, Portugal', 41.1830, -8.6890, 'TRACK002', CURRENT_DATE + INTERVAL '3 days', 'Rua António José de Almeida 20, 3040-072 Coimbra, Portugal', 40.1970, -8.4450),

  -- Order 3: Seller near Faro (Olhão) -> Logistics Center Faro -> Delivery in Coimbra (Cernache)
  (103, 503, 'Dona Alzira, Rua do Comércio 22, 8700-343 Olhão, Portugal', 37.0270, -7.8410, 'TRACK003', CURRENT_DATE + INTERVAL '7 days', 'Rua Principal 8, 3020-901 Cernache, Coimbra, Portugal', 40.2360, -8.5120),
  
  -- Order 4: Seller near Lisboa (Montijo) -> Main Warehouse -> Madeira Hub (delivery in Funchal, outskirts)
  (104, 501, 'Dona Lurdes, Rua dos Pescadores 45, 2870-108 Montijo, Portugal', 38.7060, -8.9730, 'TRACK004', CURRENT_DATE + INTERVAL '10 days', 'Estrada Monumental 390, 9000-250 Funchal, Madeira', 32.6500, -16.9300),
  
  -- Order 5: Seller near Porto (Vila Nova de Gaia) -> Distribution Center Porto -> Açores Hub (delivery in Ponta Delgada, outskirts)
  (105, 502, 'Dona Maria, Avenida da República 1200, 4430-192 Vila Nova de Gaia, Portugal', 41.1230, -8.6100, 'TRACK005', CURRENT_DATE + INTERVAL '12 days', 'Rua do Loreto 15, 9500-418 Ponta Delgada, Açores', 37.7500, -25.6700),

  -- Order 6: Seller near Évora (Vila Viçosa) -> Regional Center Évora -> Delivery in Viseu (Centro)
  (106, 504, 'Dona Beatriz, Largo do Conde 8, 7160-251 Vila Viçosa, Portugal', 38.7830, -7.4160, 'TRACK006', CURRENT_DATE + INTERVAL '8 days', 'Rua Formosa 120, 3500-161 Viseu, Portugal', 40.6566, -7.9122);

-- Insert status history for each order
INSERT INTO order_status_history (order_id, order_status, note, order_location, storage_id)
VALUES
  -- Order 1: Seller (Lisboa) -> Main Warehouse Lisboa -> Delivery
  (1, 'PROCESSING', 'Order received from seller', 'Dona Lurdes', NULL),
  (1, 'SHIPPED', 'Picked up from seller, arrived at warehouse', 'Main Warehouse Lisboa', 1),
  (1, 'IN TRANSIT', 'Package ready for delivery', 'Main Warehouse Lisboa', 1),
  
  -- Order 2: Seller (Porto) -> Distribution Center Porto -> Coimbra Hub
  (2, 'PROCESSING', 'Order received from seller', 'Dona Maria', NULL),
  (2, 'SHIPPED', 'Picked up from seller', 'Distribution Center Porto', 2),
  (2, 'IN TRANSIT', 'In transit via Coimbra', 'Regional Hub Coimbra', 3),
  
  -- Order 3: Seller (Faro) -> Logistics Center Faro (Cancelled)
  (3, 'PROCESSING', 'Order received from seller', 'Dona Alzira', NULL),
  (3, 'SHIPPED', 'Picked up from seller', 'Logistics Center Faro', 4),
  (3, 'CANCELLED', 'Customer requested cancellation', 'Logistics Center Faro', 4),
  
  -- Order 4: Seller (Lisboa) -> Main Warehouse Lisboa -> Madeira Hub
  (4, 'PROCESSING', 'Order received from seller', 'Dona Lurdes', NULL),
  (4, 'SHIPPED', 'Picked up from seller', 'Main Warehouse Lisboa', 1),
  (4, 'IN TRANSIT', 'Package at warehouse, preparing for air transit', 'Main Warehouse Lisboa', 1),
  (4, 'IN TRANSIT', 'Arrived at Lisbon Airport', 'Lisbon Airport (Cargo)', 12),
  (4, 'IN TRANSIT', 'Arrived at Funchal Airport', 'Funchal Airport (Cargo)', 14),
  (4, 'IN TRANSIT', 'Arrived at Madeira hub', 'Madeira Hub Funchal', 10),
  (4, 'OUT FOR DELIVERY', 'Out for delivery in Funchal', 'Madeira Hub Funchal', 10),
  
  -- Order 5: Seller (Porto) -> Distribution Center Porto -> Açores Hub
  (5, 'PROCESSING', 'Order received from seller', 'Dona Maria', NULL),
  (5, 'SHIPPED', 'Picked up from seller', 'Distribution Center Porto', 2),
  (5, 'IN TRANSIT', 'Package at warehouse, preparing for air transit', 'Distribution Center Porto', 2),
  (5, 'IN TRANSIT', 'Arrived at Porto Airport', 'Porto Airport (Cargo)', 13),
  (5, 'IN TRANSIT', 'Arrived at Ponta Delgada Airport', 'Ponta Delgada Airport (Cargo)', 15),
  (5, 'IN TRANSIT', 'Arrived at Açores hub', 'Açores Hub Ponta Delgada', 11),

  -- Order 6: Seller (Évora) -> Regional Center Évora -> Regional Hub Viseu (Delivered)
  (6, 'PROCESSING', 'Order received from seller', 'Dona Beatriz', NULL),
  (6, 'SHIPPED', 'Picked up from seller', 'Regional Center Évora', 6),
  (6, 'IN TRANSIT', 'In transit to Viseu', 'Regional Center Évora', 6),
  (6, 'OUT FOR DELIVERY', 'Out for delivery in Viseu', 'Regional Hub Viseu', 9),
  (6, 'DELIVERED', 'Delivered to customer in Viseu', 'Regional Hub Viseu', 9);

-- Insert products for each order
INSERT INTO order_products (order_id, product_id, quantity)
VALUES
  (1, 201, 2),
  (1, 202, 1),
  (2, 203, 3),
  (3, 204, 1),
  (3, 205, 2),
  (4, 201, 1),
  (4, 205, 3),
  (5, 202, 1),
  (5, 204, 2),
  (6, 201, 1),
  (6, 203, 2);
