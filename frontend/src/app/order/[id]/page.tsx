"use client";

import { useEffect, useState } from "react";
import { useParams } from "next/navigation";
import { OrderData, OrderProduct, OrderStatus } from "@/app/types";

export default function OrderPage() {
  const { id } = useParams();
  const [order, setOrder] = useState<OrderData | null>(null);
  const [orderHistory, setOrderHistory] = useState<OrderStatus[] | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetch(`${process.env.NEXT_PUBLIC_API_URL}/order/${id}`)
      .then(res => res.json())
      .then(data => {
        const o = data.order;
        setOrder({
          tracking_code: o.Tracking_Code,
          delivery_estimates: o.Delivery_Estimates,
          delivery_address: o.Delivery_Address,
          created_at: o.Created_At,
          price: o.Price,
          products: o.Products?.map((p: any) => ({
            product_id: p.ProductID,
            name: p.Product ? p.Product.Name : "Unknown",
            price: p.Product ? p.Product.Price : 0,
            quantity: p.Quantity
          })) || [],
          statusHistory: []
        });
      })
      .finally(() => setLoading(false));
      fetch(`${process.env.NEXT_PUBLIC_API_URL}/order/history/${id}`)
      .then(res => res.json())
      .then(data => {
        const o = data.order_status_history;
        let history: OrderStatus[] = []; 
        o.forEach((element: any) => {
          history.push({
            order_status: element.Order_Status,
            note: element.Note,
            order_location: element.Order_Location,
            timestamp: element.Timestamp_History,
            order_id: element.Order_ID,
            storage_id: element.Storage_ID,
            Storage: element.Storage ? {
              Id: element.Storage.Id,
              Name: element.Storage.Name,
              Address: element.Storage.Address,
              Latitude: element.Storage.Latitude,
              Longitude: element.Storage.Longitude,
              Created_At: element.Storage.Created_At
            } : undefined,
            order: {
              tracking_code: element.Order?.Tracking_Code,
              delivery_estimates: element.Order?.Delivery_Estimates,
              delivery_address: element.Order?.Delivery_Address,
              created_at: element.Order?.Created_At,
              price: element.Order?.Price,
              products: element.Order?.Products?.map((p: any) => ({
                product_id: p.ProductID,
                name: p.Product ? p.Product.Name : "Unknown",
                price: p.Product ? p.Product.Price : 0,
                quantity: p.Quantity
              })) || [],
              statusHistory: []
            }
          })
        });
        setOrderHistory(history);
      })
      .finally(() => setLoading(false));
  }, [id]);


  if (loading) return <p>Loading order...</p>;
  if (!order) return <p>Order not found</p>;

  return (
    <div className="max-w-4xl mx-auto p-8 space-y-8">
      {/* Map Placeholder */}
      <section className="w-full h-64 bg-gray-200 rounded-lg flex items-center justify-center">
        <p className="text-gray-500">Map will be displayed here</p>
      </section>
      {/* Header */}
      <header className="flex justify-between items-start">
        <div>
          <p className="font-bold">Order Number: {id}</p>
          <p>Order date: {new Date(order.created_at).toDateString()}</p>
          <p>Delivery: {order.delivery_address}</p>
          <p>Price: {order.price}€</p>
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
              <p className="font-semibold">{p.name}</p>
              <p>Product ID: {p.product_id}</p>
              <p>Quantity: {p.quantity}</p>
              <p>Price: {p.price}€</p>
            </li>
          ))}
        </ul>
      </section>

      {/* Tracking Steps */}
      <section>
        <h2 className="font-semibold text-lg mb-2">Tracking History</h2>
        {orderHistory && orderHistory.length > 0 ? (
          <ul className="space-y-4">
            {orderHistory.map((s, idx) => (
              <li key={idx} className="flex items-start gap-4">
                <div className="w-6 h-6 rounded-full bg-blue-600 shrink-0" />
                <div>
                  <p className="font-bold">{s.order_status}</p>
                  <p>{s.note} ({s.order_location})</p>
                  {s.Storage && (
                    <p className="text-sm text-gray-600">
                      Storage: {s.Storage.Name} - {s.Storage.Address}
                    </p>
                  )}
                </div>
              </li>
            ))}
          </ul>
        ) : (
          <p className="text-gray-500 italic">No tracking history available.</p>
        )}
      </section>
    </div>
  );
}
