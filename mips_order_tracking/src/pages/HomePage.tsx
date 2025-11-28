import { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';
import type { OrderData, BackendOrder, BackendOrderProduct, BackendOrderStatus } from '../types';

export default function OrdersPage() {
    const [orders, setOrders] = useState<OrderData[] | null>(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const [order_by, setOrderBy] = useState<string>("newest");
    const [statusFilter, setStatusFilter] = useState<string>("all");

    const apiUrl = process.env.PUBLIC_API_URL || 'http://localhost:8080';

    useEffect(() => {
        setLoading(true);
        handleOrders();
    }, [order_by, statusFilter]);

    function handleOrders(){
        fetch(`${apiUrl}/orders?order_by=${order_by}`)
            .then((res) => {
                if (!res.ok) setError("Could not load the orders");
                return res.json();
            })
            .then((data) => {
                console.log(data)
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
                    statusHistory:ord.Order_Status?.map((p: BackendOrderStatus) => ({
                        order_status: p.Order_Status,
                        note: p.Note,
                        order_location: p.Order_Location,
                        timestamp: new Date(p.Timestamp_History),
                        order_id: p.Order_ID,
                        order: ord ? {
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
                            products: [],
                            statusHistory: []
                        } : null,
                        storage_id: p.Storage_ID
                    })) || [],
                }));
                console.log(parsed)
                const filtered = parsed.filter((o) => {
                    const matchesStatus =
                        statusFilter === "all" || o.statusHistory?.[0]?.order_status === statusFilter;
                    return matchesStatus;
                });
                setOrders(filtered);
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
                 {/* Status filter */}
                <label className="flex flex-col text-sm font-medium text-gray-700 dark:text-gray-300">
                    Status
                    <select
                    value={statusFilter}
                    onChange={(e) => setStatusFilter(e.target.value)}
                    className="mt-1 block w-48 rounded-xl border border-gray-300 dark:border-gray-700 
                                bg-white dark:bg-gray-900 text-gray-900 dark:text-gray-100 
                                px-3 py-2 shadow-sm focus:outline-none focus:ring-2 
                                focus:ring-indigo-500 dark:focus:ring-indigo-400 
                                transition-colors"
                    >
                    <option value="all">All statuses</option>
                    <option value="PROCESSING">Processing</option>
                    <option value="SHIPPED">Shipped</option>
                    <option value="IN TRANSIT">In Transit</option>
                    <option value="OUT FOR DELIVERY">Out for Delivery</option>
                    <option value="DELIVERED">Delivered</option>
                    <option value="CANCELLED">Cancelled</option>
                    <option value="RETURNED">Returned</option>
                    <option value="FAILED DELIVERY">Failed Delivery</option>
                    </select>
                </label>

                {/* Sort selector */}
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