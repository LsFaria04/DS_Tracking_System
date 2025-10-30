interface OrderStatus {
  order_status: string;
  note: string;
  order_location: string;
  timestamp: Date;
  order: OrderData | null;
  order_id: number;
}

interface OrderProduct {
  product_id: number;
  name: string;
  price: number;
  quantity: number;
}

interface OrderData {
  tracking_code: string;
  delivery_estimates: string;
  delivery_address: string;
  created_at: string;
  price: string;
  products: OrderProduct[];
  statusHistory: OrderStatus[];
}

export type {OrderData, OrderProduct, OrderStatus}