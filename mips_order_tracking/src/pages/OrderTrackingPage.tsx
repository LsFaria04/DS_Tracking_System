import { useEffect, useState, lazy, Suspense } from "react";
import { useParams } from "react-router-dom";
import type { BackendOrder, BackendOrderProduct, BackendOrderStatus, OrderData, OrderStatus, VerificationResult } from "../types";

const OrderMap = lazy(() => import('../components/OrderMap'));

export default function OrderPage() {
    const { id } = useParams<{ id: string }>(); 
    const [order, setOrder] = useState<OrderData | null>(null);
    const [orderHistory, setOrderHistory] = useState<OrderStatus[] | null>(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const [verifying, setVerifying] = useState(false);
    const [verificationResult, setVerificationResult] = useState<VerificationResult | null>(null);

    useEffect(() => {
        const apiUrl = process.env.PUBLIC_API_URL || 'http://localhost:8080';
        
        fetch(`${apiUrl}/order/${id}`)
            .then(res => {
                if (!res.ok) throw new Error('Failed to fetch order');
                return res.json();
            })
            .then(data => {
                const o = data.order as BackendOrder;
                console.log(o)
                setOrder({
                    tracking_code: o.Tracking_Code,
                    delivery_estimate: o.Delivery_Estimate,
                    delivery_address: o.Delivery_Address,
                    delivery_latitude: o.Delivery_Latitude,
                    delivery_longitude: o.Delivery_Longitude,
                    seller_address: o.Seller_Address,
                    seller_latitude: o.Seller_Latitude,
                    seller_longitude: o.Seller_Longitude,
                    created_at: o.Created_At,
                    price: o.Price.toString(),
                    products: o.Products?.map((p: BackendOrderProduct) => ({
                        product_id: p.Product_ID,
                        name: p.Product_Name_At_Purchase,
                        price: p.Product_Price_At_Purchase,
                        quantity: p.Quantity
                    })) || [],
                    statusHistory: []
                });
            })
            .catch(() => {
                setError('Failed to load order. Please try again.');
            })
            .finally(() => setLoading(false));
            
            fetch(`${apiUrl}/order/history/${id}`)
            .then(res => {
                if (!res.ok) throw new Error('Failed to fetch order history');
                return res.json();
            })
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
                            delivery_estimate: element.Order.Delivery_Estimate,
                            delivery_address: element.Order.Delivery_Address,
                            delivery_latitude: element.Order.Delivery_Latitude,
                            delivery_longitude: element.Order.Delivery_Longitude,
                            seller_address: element.Order.Seller_Address,
                            seller_latitude: element.Order.Seller_Latitude,
                            seller_longitude: element.Order.Seller_Longitude,
                            created_at: element.Order.Created_At,
                            price: element.Order.Price.toString(),
                            products: element.Order.Products?.map((p: BackendOrderProduct) => ({
                                product_id: p.Product_ID,
                                name: p.Product_Name_At_Purchase,
                                price: p.Product_Price_At_Purchase,
                                quantity: p.Quantity
                            })) || [],
                            statusHistory: []
                        } : null
                    })
                });
                setOrderHistory(history);
            })
            .catch(() => {
                // History is optional, don't set error if it fails
                console.warn('Failed to load order history');
            });
    }, [id]);

    const handleVerifyBlockchain = async () => {
        setVerifying(true);
        setVerificationResult(null);
        
        try {
            const apiUrl = process.env.PUBLIC_API_URL || 'http://localhost:8080';
            const response = await fetch(`${apiUrl}/order/verify/${id}`);
            const data = await response.json();
            
            if (response.ok) {
                setVerificationResult(data);
            } else {
                setVerificationResult({
                    verified: false,
                    total_updates: 0,
                    verified_updates: 0,
                    blockchain_hashes: 0,
                    status: 'ERROR',
                    message: data.error || 'Failed to verify order',
                });
            }
        } catch {
            setVerificationResult({
                verified: false,
                total_updates: 0,
                verified_updates: 0,
                blockchain_hashes: 0,
                status: 'ERROR',
                message: 'Network error while verifying',
            });
        } finally {
            setVerifying(false);
        }
    };



    if (loading) return (
        <div className="max-w-5xl mx-auto p-6 md:p-8 space-y-6">
            {/* Map skeleton */}
            <div className="w-full h-96 bg-gray-200 dark:bg-gray-800 rounded-2xl animate-pulse"></div>
            
            {/* Header skeleton */}
            <div className="bg-white dark:bg-gray-900 border border-gray-200 dark:border-gray-800 rounded-2xl p-6">
                <div className="flex flex-col md:flex-row md:justify-between gap-6">
                    <div className="flex-1 space-y-4">
                        <div className="h-8 bg-gray-200 dark:bg-gray-800 rounded-lg w-48 animate-pulse"></div>
                        <div className="space-y-2">
                            <div className="h-4 bg-gray-200 dark:bg-gray-800 rounded w-64 animate-pulse"></div>
                            <div className="h-4 bg-gray-200 dark:bg-gray-800 rounded w-64 animate-pulse"></div>
                            <div className="h-4 bg-gray-200 dark:bg-gray-800 rounded w-96 animate-pulse"></div>
                        </div>
                    </div>
                    <div className="flex flex-col gap-4 min-w-[200px]">
                        <div className="h-16 bg-gray-200 dark:bg-gray-800 rounded-lg animate-pulse"></div>
                        <div className="h-20 bg-gray-200 dark:bg-gray-800 rounded-lg animate-pulse"></div>
                    </div>
                </div>
            </div>
            
            {/* Products skeleton */}
            <div>
                <div className="h-6 bg-gray-200 dark:bg-gray-800 rounded w-32 mb-6 animate-pulse"></div>
                <div className="space-y-3">
                    {[1, 2].map(i => (
                        <div key={i} className="bg-white dark:bg-gray-900 border border-gray-200 dark:border-gray-800 rounded-2xl p-5">
                            <div className="h-5 bg-gray-200 dark:bg-gray-800 rounded w-48 mb-4 animate-pulse"></div>
                            <div className="h-4 bg-gray-200 dark:bg-gray-800 rounded w-32 animate-pulse"></div>
                        </div>
                    ))}
                </div>
            </div>
            
            {/* Tracking history skeleton */}
            <div>
                <div className="h-6 bg-gray-200 dark:bg-gray-800 rounded w-40 mb-6 animate-pulse"></div>
                <div className="space-y-0">
                    {[1, 2, 3].map(i => (
                        <div key={i} className="relative flex pb-6">
                            <div className="relative flex flex-col items-center w-12 shrink-0">
                                <div className="w-4 h-4 bg-gray-300 dark:bg-gray-700 rounded-full animate-pulse"></div>
                                {i < 3 && <div className="w-0.5 flex-1 bg-gray-300 dark:bg-gray-700"></div>}
                            </div>
                            <div className="flex-1">
                                <div className="bg-white dark:bg-gray-900 border border-gray-200 dark:border-gray-800 rounded-2xl p-5">
                                    <div className="h-5 bg-gray-200 dark:bg-gray-800 rounded w-32 mb-3 animate-pulse"></div>
                                    <div className="h-4 bg-gray-200 dark:bg-gray-800 rounded w-64 mb-2 animate-pulse"></div>
                                    <div className="h-4 bg-gray-200 dark:bg-gray-800 rounded w-48 animate-pulse"></div>
                                </div>
                            </div>
                        </div>
                    ))}
                </div>
            </div>
        </div>
    );
    
    if (error) return (
        <div className="min-h-screen flex flex-col items-center justify-center p-6">
            <div className="text-center max-w-md">
                <p className="text-6xl mb-4"></p>
                <p className="text-xl font-semibold text-gray-900 dark:text-white mb-2">Unable to load order</p>
                <p className="text-gray-500 dark:text-gray-400 mb-6">{error}</p>
                <button
                    onClick={() => window.location.reload()}
                    className="px-6 py-3 bg-blue-500 hover:bg-blue-600 text-white font-medium rounded-lg transition-colors"
                >
                    Try Again
                </button>
            </div>
        </div>
    );
    
    if (!order) return (
        <div className="min-h-screen flex flex-col items-center justify-center">
            <div className="text-center">
                <p className="text-6xl mb-4"></p>
                <p className="text-xl font-semibold text-gray-900 dark:text-white mb-2">Order not found</p>
                <p className="text-gray-500 dark:text-gray-400">The order you&apos;re looking for doesn&apos;t exist or has been removed.</p>
            </div>
        </div>
    );

    return (
        <div className="max-w-5xl mx-auto p-6 md:p-8 space-y-6">
            {/* Map */}
            <section>
                {/* 5. Wrapped lazy component in <Suspense> */}
                <Suspense fallback={
                    <div className="w-full h-96 bg-gray-200 rounded-lg flex items-center justify-center">
                        <p className="text-gray-500">Loading map...</p>
                    </div>
                }>
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
                </Suspense>
            </section>

            {/* Header */}
            <header className="bg-white dark:bg-gray-900 border border-gray-200 dark:border-gray-800 rounded-2xl p-6">
                <div className="flex flex-col md:flex-row md:justify-between md:items-start gap-6">
                    <div className="flex-1 space-y-4">
                        <h1 className="text-3xl font-semibold text-gray-900 dark:text-white">Order {id}</h1>
                        <div className="space-y-2 text-sm">
                            <div className="flex items-center gap-3">
                                <span className="text-gray-500 dark:text-gray-400 w-28">Ordered</span>
                                <span className="text-gray-900 dark:text-white">{new Date(order.created_at).toLocaleDateString()}</span>
                            </div>
                            <div className="flex items-center gap-3">
                                <span className="text-gray-500 dark:text-gray-400 w-28">Est. Delivery</span>
                                <span className="text-gray-900 dark:text-white">{new Date(order.delivery_estimate).toLocaleDateString()}</span>
                            </div>
                            <div className="flex items-start gap-3">
                                <span className="text-gray-500 dark:text-gray-400 w-28">Delivery To</span>
                                <span className="text-gray-900 dark:text-white">{order.delivery_address}</span>
                            </div>
                        </div>
                    </div>
                    <div className="flex flex-col items-start md:items-end gap-4 min-w-[200px]">
                        <div className="text-right w-full">
                            <p className="text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider mb-1">Tracking Code</p>
                            <p className="text-base font-mono font-medium text-gray-900 dark:text-white">{order.tracking_code}</p>
                        </div>
                        <div className="text-right w-full pt-4 border-t border-gray-200 dark:border-gray-800">
                            <p className="text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider mb-1">Total</p>
                            <p className="text-3xl font-semibold text-gray-900 dark:text-white">{order.price}€</p>
                        </div>
                        
                        {/* Blockchain Verification */}
                        <div className="w-full pt-4 border-t border-gray-200 dark:border-gray-800">
                            <button
                                onClick={handleVerifyBlockchain}
                                disabled={verifying}
                                className="w-full px-4 py-2 bg-blue-500 hover:bg-blue-600 disabled:bg-gray-400 text-white text-sm font-medium rounded-lg transition-colors disabled:cursor-not-allowed flex items-center justify-center gap-2"
                            >
                                {verifying ? (
                                    <>
                                        <span className="inline-block w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin"></span>
                                        Verifying...
                                    </>
                                ) : (
                                    'Verify on Blockchain'
                                )}
                            </button>
                            
                            {verificationResult && (
                                <div className={`mt-3 p-3 rounded-lg text-sm ${
                                    verificationResult.status === 'VERIFIED' 
                                        ? 'bg-green-50 dark:bg-green-900/20 border border-green-200 dark:border-green-800'
                                        : verificationResult.status === 'ERROR'
                                        ? 'bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800'
                                        : 'bg-orange-50 dark:bg-orange-900/20 border border-orange-200 dark:border-orange-800'
                                }`}>
                                    <div className="flex items-center gap-2 mb-1">
                                        <span className={`font-semibold ${
                                            verificationResult.status === 'VERIFIED'
                                                ? 'text-green-700 dark:text-green-300'
                                                : verificationResult.status === 'ERROR'
                                                ? 'text-red-700 dark:text-red-300'
                                                : 'text-orange-700 dark:text-orange-300'
                                        }`}>
                                            {verificationResult.status.replace('_', ' ')}
                                        </span>
                                    </div>
                                    <p className={`text-xs ${
                                        verificationResult.status === 'VERIFIED'
                                            ? 'text-green-600 dark:text-green-400'
                                            : verificationResult.status === 'ERROR'
                                            ? 'text-red-600 dark:text-red-400'
                                            : 'text-orange-600 dark:text-orange-400'
                                    }`}>
                                        {verificationResult.message}
                                    </p>
                                    {verificationResult.status !== 'ERROR' && (
                                        <>
                                            <p className="text-xs text-gray-600 dark:text-gray-400 mt-2">
                                                {verificationResult.verified_updates}/{verificationResult.total_updates} updates verified
                                            </p>
                                            {verificationResult.contract_address && (
                                                <p className="text-xs text-gray-600 dark:text-gray-400 mt-1">
                                                    <span className="font-medium">Contract:</span> {verificationResult.contract_address}
                                                </p>
                                            )}
                                            {verificationResult.transaction_hashes && verificationResult.transaction_hashes.length > 0 && (
                                                <p className="text-xs text-gray-600 dark:text-gray-400 mt-2">
                                                    <span className="font-medium">✓</span> Blockchain links shown in tracking history below
                                                </p>
                                            )}
                                        </>
                                    )}
                                </div>
                            )}
                        </div>
                    </div>
                </div>
            </header>

            {/* Products */}
            <section>
                <h2 className="text-2xl font-semibold mb-6 text-gray-900 dark:text-white">Order Items</h2>
                {order.products && order.products.length > 0 ? (
                    <div className="space-y-3">
                        {order.products.map((p, idx) => (
                            <div key={p.product_id || idx} className="bg-white dark:bg-gray-900 border border-gray-200 dark:border-gray-800 rounded-2xl p-5 hover:border-gray-300 dark:hover:border-gray-700 transition-colors">
                                <div className="flex flex-col sm:flex-row justify-between items-start gap-4">
                                    <div className="flex-1">
                                        <h3 className="text-lg font-medium text-gray-900 dark:text-white mb-1">{p.name}</h3>
                                        <p className="text-sm text-gray-500 dark:text-gray-400">
                                            Product ID: {p.product_id}
                                        </p>
                                    </div>
                                    <div className="text-right">
                                        <p className="text-xl font-semibold text-gray-900 dark:text-white">{p.price}€</p>
                                        <p className="text-xs text-gray-500 dark:text-gray-400">per unit</p>
                                    </div>
                                </div>
                                
                                <div className="mt-4 flex justify-between items-center pt-4 border-t border-gray-200 dark:border-gray-800">
                                    <div className="flex items-center gap-2">
                                        <span className="text-sm text-gray-500 dark:text-gray-400">Quantity</span>
                                        <span className="text-sm font-medium text-gray-900 dark:text-white">
                                            {p.quantity}
                                        </span>
                                    </div>
                                    <div className="text-right">
                                        <span className="text-sm text-gray-500 dark:text-gray-400 mr-2">Subtotal</span>
                                        <span className="text-lg font-semibold text-gray-900 dark:text-white">{(p.price * p.quantity).toFixed(2)}€</span>
                                    </div>
                                </div>
                            </div>
                        ))}
                        
                        {/* Total summary */}
                        <div className="mt-6 bg-gray-50 dark:bg-gray-900 border border-gray-200 dark:border-gray-800 rounded-2xl p-5">
                            <div className="flex justify-between items-center">
                                <div>
                                    <p className="text-sm font-medium text-gray-900 dark:text-white">Total</p>
                                    <p className="text-xs text-gray-500 dark:text-gray-400">{order.products.length} item{order.products.length !== 1 ? 's' : ''}</p>
                                </div>
                                <p className="text-3xl font-semibold text-gray-900 dark:text-white">
                                    {order.price}€
                                </p>
                            </div>
                        </div>
                    </div>
                ) : (
                    <div className="text-center py-12 bg-gray-50 dark:bg-gray-900 border border-gray-200 dark:border-gray-800 rounded-2xl">
                        <p className="text-gray-500 dark:text-gray-400">No products in this order.</p>
                    </div>
                )}
            </section>

            {/* Tracking Steps */}
            <section>
                <h2 className="text-2xl font-semibold mb-6 text-gray-900 dark:text-white">Tracking History</h2>
                {orderHistory && orderHistory.length > 0 ? (
                    <ul className="space-y-0">
                        {orderHistory.map((s, idx) => {
                            const isFirst = idx === 0;
                            const isLast = idx === orderHistory.length - 1;
                            
                            // Define status colors
                            const statusInfo = {
                                'PROCESSING': { color: 'bg-blue-500', border: 'border-blue-500' },
                                'SHIPPED': { color: 'bg-blue-500', border: 'border-blue-5POST-t-t' },
                                'IN TRANSIT': { color: 'bg-blue-500', border: 'border-blue-500' },
                                'OUT FOR DELIVERY': { color: 'bg-blue-500', border: 'border-blue-500' },
                                'DELIVERED': { color: 'bg-green-500', border: 'border-green-500' },
                                'CANCELLED': { color: 'bg-red-500', border: 'border-red-500' },
                                'RETURNED': { color: 'bg-orange-500', border: 'border-orange-500' },
                                'FAILED DELIVERY': { color: 'bg-red-500', border: 'border-red-500' }
                            };
                            
                            const status = statusInfo[s.order_status as keyof typeof statusInfo] || { 
                                color: 'bg-gray-400',
                                border: 'border-gray-400'
                            };
                            
                            // Get the blockchain transaction hash for this status update
                            // Transaction hashes align with the order returned by the verification API
                            // (same index as the orderHistory entries shown), so use the same index
                            const txHash = verificationResult?.transaction_hashes?.[idx];
                            
                            return (
                                <li key={idx} className="relative flex">
                                    {/* Left timeline column */}
                                    <div className="relative flex flex-col items-center w-12 shrink-0">
                                        {/* Top connecting line - only if not first */}
                                        {!isFirst && (
                                            <div className="w-0.5 h-6 bg-gray-300 dark:bg-gray-700"></div>
                                        )}
                                        
                                        {/* The dot - larger for the last (current state) one */}
                                        <div className={`${status.color} rounded-full ${isLast ? 'w-8 h-8' : 'w-6 h-6'} shrink-0 border-4 border-white dark:border-gray-950`}></div>
                                        
                                        {/* Bottom connecting line - only if not last */}
                                        {!isLast && (
                                            <div className="w-0.5 flex-1 bg-gray-300 dark:bg-gray-700"></div>
                                        )}
                                    </div>
                                    
                                    {/* Content card */}
                                    <div className={`flex-1 pb-6 ${isLast ? 'pb-0' : ''}`}>
                                        <div className="bg-white dark:bg-gray-900 p-5 rounded-2xl border border-gray-200 dark:border-gray-800 hover:border-gray-300 dark:hover:border-gray-700 transition-colors">
                                            <div className="flex justify-between items-start mb-3">
                                                <h3 className="font-semibold text-base text-gray-900 dark:text-white">{s.order_status}</h3>
                                                <span className="text-xs text-gray-500 dark:text-gray-400 whitespace-nowrap ml-4">
                                                    {new Date(s.timestamp).toLocaleString()}
                                                </span>
                                            </div>
                                            
                                            <p className="text-sm text-gray-600 dark:text-gray-400 mb-2">{s.note}</p>
                                            
                                            <div className="text-sm text-gray-500 dark:text-gray-400">
                                                {s.order_location}
                                            </div>
                                            
                                            {s.Storage && (
                                                <div className="mt-3 pt-3 border-t border-gray-200 dark:border-gray-800">
                                                    <div className="text-sm">
                                                        <p className="font-medium text-gray-900 dark:text-white">{s.Storage.Name}</p>
                                                        <p className="text-gray-500 dark:text-gray-400 text-xs mt-1">{s.Storage.Address}</p>
                                                    </div>
                                                </div>
                                            )}
                                            
                                            {/* Blockchain verification link */}
                                            {txHash && (
                                                <div className="mt-3 pt-3 border-t border-gray-200 dark:border-gray-800">
                                                    <div className="flex items-center gap-2">
                                                        <span className="text-xs font-medium text-green-600 dark:text-green-400">Verified on Blockchain</span>
                                                    </div>
                                                    <a
                                                        href={`https://sepolia.etherscan.io/tx/${txHash}`}
                                                        target="_blank"
                                                        rel="noopener noreferrer"
                                                        className="text-xs text-blue-600 dark:text-blue-400 hover:underline font-mono break-all mt-1 block"
                                                        title="View transaction on Etherscan"
                                                    >
                                                        {txHash}
                                                    </a>
                                                </div>
                                            )}
                                        </div>
                                    </div>
                                </li>
                            );
                        })}
                    </ul>
                ) : (
                    <div className="text-center py-12 bg-gray-50 dark:bg-gray-900 border border-gray-200 dark:border-gray-800 rounded-2xl">
                        <p className="text-gray-500 dark:text-gray-400">No tracking history available yet.</p>
                    </div>
                )}
            </section>
        </div>
    );
}