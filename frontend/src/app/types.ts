interface OrderProduct {
  product_id: number;
  name: string;
  price: number;
  quantity: number;
}

interface Storage {
  Id: number;
  Name: string;
  Address: string;
  Latitude: number;
  Longitude: number;
  Created_At: string;
}

interface OrderStatus {
  order_status: string;
  note: string;
  order_location: string;
  timestamp: Date;
  order: OrderData | null;
  order_id: number;
  storage_id?: number;
  Storage?: Storage;
}

interface OrderData {
  tracking_code: string;
  delivery_estimates: string;
  delivery_address: string;
  delivery_latitude?: number;
  delivery_longitude?: number;
  seller_address?: string;
  seller_latitude?: number;
  seller_longitude?: number;
  created_at: string;
  price: string;
  products: OrderProduct[];
  statusHistory: OrderStatus[];
}

export type {OrderData, OrderProduct, OrderStatus, Storage}