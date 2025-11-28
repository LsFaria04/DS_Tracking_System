import { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';
import type { OrderData, BackendOrder, BackendOrderProduct } from '../types';

export default function OrdersPage() {
    const [orders, setOrders] = useState<OrderData[] | null>(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const [order_by, setOrderBy] = useState<string>("newest");
    const apiUrl = process.env.PUBLIC_API_URL || 'http://localhost:8080';

    useEffect(() => {
        setLoading(true);
        handleOrders();
    }, [order_by]);

    function handleOrders(){
        console.log(order_by)
        fetch(`${apiUrl}/orders?order_by=${order_by}`)
            .then((res) => {
                if (!res.ok) throw new Error('Failed to fetch orders');
                return res.json();
            })
            .then((data) => {
                const o = data.orders as BackendOrder[];
                const parsed = o.map((ord) => ({
                    id: ord.Id,
                    tracking_code: ord.Tracking_Code,
                    delivery_estimate: ord.Delivery_Estimate,
                    delivery_address: ord.Delivery_Address,
                    delivery_latitude: ord.Delivery_Latitude,
                    delivery_longitude: ord.Delivery_Longitude,
                    seller_address: ord.Seller_Address,
                    seller_latitude: ord.Seller_Latitude,
                    seller_longitude: ord.Seller_Longitude,
                    created_at: ord.Created_At,
                    price: ord.Price.toString(),
                    products: ord.Products?.map((p: BackendOrderProduct) => ({
                        product_id: p.Product_ID,
                        name: p.Product_Name_At_Purchase,
                        price: p.Product_Price_At_Purchase,
                        quantity: p.Quantity,
                    })) || [],
                    statusHistory: [],
                }));
                setOrders(parsed);
            })
            .catch((err) => {
                console.warn(err);
                setError('Failed to load orders');
            })
            .finally(() => setLoading(false));
    }
    
    if (loading) return (
        <div className="p-6">
            <div className="h-8 bg-gray-200 dark:bg-gray-800 rounded w-48 mb-4 animate-pulse"></div>
            <div className="space-y-3">
                {[1,2,3].map(i => (
                    <div key={i} className="bg-white dark:bg-gray-900 border rounded p-4 animate-pulse"></div>
                ))}
            </div>
        </div>
    );

    if (error) return (
        <div className="min-h-screen flex items-center justify-center">
            <div className="text-center">
                <p className="text-xl font-semibold mb-2">Unable to load orders</p>
                <p className="text-sm text-gray-500 mb-6">{error}</p>
            </div>
        </div>
    );

    return (
        <div className="max-w-5xl mx-auto p-6 md:p-8 space-y-6">
            <h1 className="text-3xl font-semibold text-gray-900 dark:text-white">Orders</h1>
            <div className="flex gap-4 mb-6">
            <label className="flex flex-col text-sm font-medium text-gray-700 dark:text-gray-300">
                Sort by:
                <select
                value={order_by}
                onChange={(e) => setOrderBy(e.target.value)}
                className="mt-1 block w-48 rounded-xl border border-gray-300 dark:border-gray-700 
                            bg-white dark:bg-gray-900 text-gray-900 dark:text-gray-100 
                            px-3 py-2 shadow-sm focus:outline-none focus:ring-2 
                            focus:ring-indigo-500 dark:focus:ring-indigo-400 
                            transition-colors"
                >
                <option value="newest">Newest first</option>
                <option value="oldest">Oldest first</option>
                </select>
            </label>
            </div>
            {orders && orders.length > 0 ? (
                <div className="space-y-3">
                {orders.map((o) => (
                    <Link to={`/order/${o.id}`} key={o.id} className="block bg-white dark:bg-gray-900 border rounded-2xl p-5 hover:border-gray-300 dark:hover:border-gray-700 transition-colors">
                        <div className="flex justify-between items-start">
                            <div>
                                <p className="font-medium text-gray-900 dark:text-white">{o.tracking_code}</p>
                                <p className="text-sm text-gray-500 dark:text-gray-400">{o.delivery_address}</p>
                            </div>
                            <div className="text-right">
                                <p className="font-semibold text-gray-900 dark:text-white">{o.price}€</p>
                                <p className="text-xs text-gray-500 dark:text-gray-400">{o.products?.length ?? 0} item{o.products?.length !== 1 ? 's' : ''}</p>
                            </div>
                        </div>
                    </Link>
                ))}
                </div>
            ) : (
                <div className="text-center py-12 bg-gray-50 dark:bg-gray-900 border rounded-2xl">
                    <p className="text-gray-500 dark:text-gray-400">No orders found.</p>
                </div>
            )}
        </div>
    );
}