"use client";

import { useEffect, useState } from "react";
import { useParams } from "next/navigation";

interface OrderStatus {
  order_status: string;
  note: string;
  order_location: string;
}

interface OrderProduct {
  product_id: number;
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

export default function OrderPage() {
  const { id } = useParams();
  const [order, setOrder] = useState<OrderData | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetch(`http://localhost:8080/order/${id}`)
      .then(res => res.json())
      .then(data => {
        const o = data.order;
        setOrder({
          tracking_code: o.Tracking_Code,
          delivery_estimates: o.Delivery_Estimates,
          delivery_address: o.Delivery_Address,
          created_at: o.Created_At,
          price: o.Price,
          products: [],
          statusHistory: [] 
        });
      })
      .finally(() => setLoading(false));
  }, [id]);

  if (loading) return <p>Loading order...</p>;
  if (!order) return <p>Order not found</p>;

  return (
    <div className="max-w-4xl mx-auto p-8 space-y-8">
      {/* Header */}
      <header className="flex justify-between items-start">
        <div>
          <p className="font-bold">Order Number: {id}</p>
          <p>Order date: {new Date(order.created_at).toDateString()}</p>
          <p>Delivery: {order.delivery_address}</p>
          <p>Price: {order.price}</p> 
          <p>Estimated Delivery: {new Date(order.delivery_estimates).toDateString()}</p>
        </div>
        <div>
          <p className="font-bold">Tracking Code:</p>
          <p>{order.tracking_code}</p>
        </div>
      </header>

      {/* Products */}
      <section>
        <h2 className="font-semibold text-lg mb-2">Products</h2>
        <ul className="space-y-2">
          {order.products.map(p => (
            <li key={p.product_id} className="border p-2 rounded">
              Product {p.product_id} — Quantity: {p.quantity}
            </li>
          ))}
        </ul>
      </section>

      {/* Tracking Steps */}
      <section>
        <h2 className="font-semibold text-lg mb-2">Tracking History</h2>
        <ul className="space-y-4">
          {order.statusHistory.map((s, idx) => (
            <li key={idx} className="flex items-start gap-4">
              <div className="w-6 h-6 rounded-full bg-blue-600 flex-shrink-0" />
              <div>
                <p className="font-bold">{s.order_status}</p>
                <p>{s.note} ({s.order_location})</p>
              </div>
            </li>
          ))}
        </ul>
      </section>
    </div>
  );
}
