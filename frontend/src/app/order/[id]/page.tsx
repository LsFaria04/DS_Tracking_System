"use client";

import { useEffect, useState } from "react";
import { useParams } from "next/navigation";
import { OrderData, OrderStatus } from "@/app/types";
import dynamic from 'next/dynamic';

// Import map component dynamically to avoid SSR issues
const OrderMap = dynamic(() => import('@/app/components/OrderMap'), {
  ssr: false,
  loading: () => (
    <div className="w-full h-96 bg-gray-200 rounded-lg flex items-center justify-center">
      <p className="text-gray-500">Loading map...</p>
    </div>
  )
});

// Backend API response types
interface BackendProduct {
  ID: number;
  Name: string;
  Price: number;
}

interface BackendOrderProduct {
  ProductID: number;
  Quantity: number;
  Product?: BackendProduct;
}

interface BackendStorage {
  Id: number;
  Name: string;
  Address: string;
  Latitude: number;
  Longitude: number;
  Created_At: string;
}

interface BackendOrder {
  Tracking_Code: string;
  Delivery_Estimates: string;
  Delivery_Address: string;
  Delivery_Latitude?: number;
  Delivery_Longitude?: number;
  Seller_Address?: string;
  Seller_Latitude?: number;
  Seller_Longitude?: number;
  Created_At: string;
  Price: number;
  Products?: BackendOrderProduct[];
}

interface BackendOrderStatus {
  Order_Status: string;
  Note: string;
  Order_Location: string;
  Timestamp_History: string;
  Order_ID: number;
  Storage_ID?: number;
  Storage?: BackendStorage;
  Order?: BackendOrder;
}

export default function OrderPage() {
  const { id } = useParams();
  const [order, setOrder] = useState<OrderData | null>(null);
  const [orderHistory, setOrderHistory] = useState<OrderStatus[] | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetch(`${process.env.NEXT_PUBLIC_API_URL}/order/${id}`)
      .then(res => res.json())
      .then(data => {
        const o = data.order as BackendOrder;
        setOrder({
          tracking_code: o.Tracking_Code,
          delivery_estimates: o.Delivery_Estimates,
          delivery_address: o.Delivery_Address,
          delivery_latitude: o.Delivery_Latitude,
          delivery_longitude: o.Delivery_Longitude,
          seller_address: o.Seller_Address,
          seller_latitude: o.Seller_Latitude,
          seller_longitude: o.Seller_Longitude,
          created_at: o.Created_At,
          price: o.Price.toString(),
          products: o.Products?.map((p: BackendOrderProduct) => ({
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
        const o = data.order_status_history as BackendOrderStatus[];
        const history: OrderStatus[] = []; 
        o.forEach((element: BackendOrderStatus) => {
          history.push({
            order_status: element.Order_Status,
            note: element.Note,
            order_location: element.Order_Location,
            timestamp: new Date(element.Timestamp_History),
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
            order: element.Order ? {
              tracking_code: element.Order.Tracking_Code,
              delivery_estimates: element.Order.Delivery_Estimates,
              delivery_address: element.Order.Delivery_Address,
              delivery_latitude: element.Order.Delivery_Latitude,
              delivery_longitude: element.Order.Delivery_Longitude,
              seller_address: element.Order.Seller_Address,
              seller_latitude: element.Order.Seller_Latitude,
              seller_longitude: element.Order.Seller_Longitude,
              created_at: element.Order.Created_At,
              price: element.Order.Price.toString(),
              products: element.Order.Products?.map((p: BackendOrderProduct) => ({
                product_id: p.ProductID,
                name: p.Product ? p.Product.Name : "Unknown",
                price: p.Product ? p.Product.Price : 0,
                quantity: p.Quantity
              })) || [],
              statusHistory: []
            } : null
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
      {/* Map */}
      <section>
        {orderHistory && (
          <OrderMap
            orderHistory={orderHistory}
            deliveryAddress={order.delivery_address}
            deliveryLatitude={order.delivery_latitude}
            deliveryLongitude={order.delivery_longitude}
            sellerAddress={order.seller_address}
            sellerLatitude={order.seller_latitude}
            sellerLongitude={order.seller_longitude}
          />
        )}
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
        <h2 className="font-semibold text-lg mb-4">Order Items</h2>
        {order.products && order.products.length > 0 ? (
          <div className="space-y-3">
            {order.products.map((p, idx) => (
              <div key={p.product_id || idx} className="border border-gray-300 p-4 rounded-lg shadow-sm hover:shadow-md transition-shadow bg-white dark:bg-gray-800 dark:border-gray-700">
                <div className="flex justify-between items-start">
                  <div className="flex-1">
                    <p className="font-semibold text-lg">{p.name}</p>
                    <p className="text-sm text-gray-600 dark:text-gray-400">Product ID: #{p.product_id}</p>
                  </div>
                  <div className="text-right">
                    <p className="font-bold text-blue-600 dark:text-blue-400">{p.price}€</p>
                    <p className="text-sm text-gray-600 dark:text-gray-400">per unit</p>
                  </div>
                </div>
                <div className="mt-3 flex justify-between items-center pt-3 border-t border-gray-200 dark:border-gray-700">
                  <p className="text-gray-700 dark:text-gray-300">Quantity: <span className="font-semibold">{p.quantity}</span></p>
                  <p className="font-semibold text-gray-900 dark:text-gray-100">Subtotal: {(p.price * p.quantity).toFixed(2)}€</p>
                </div>
              </div>
            ))}
            <div className="mt-4 pt-4 border-t-2 border-gray-300 dark:border-gray-700">
              <div className="flex justify-between items-center">
                <p className="text-xl font-bold">Total Order Value:</p>
                <p className="text-2xl font-bold text-blue-600 dark:text-blue-400">{order.price}€</p>
              </div>
            </div>
          </div>
        ) : (
          <p className="text-gray-500 italic">No products in this order.</p>
        )}
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
