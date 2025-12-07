import { useEffect, useState, Suspense, useCallback } from "react";
import { useParams } from "react-router-dom";
import type { BackendOrder, BackendOrderProduct, BackendOrderStatus, OrderData, OrderStatus, VerificationResult } from "../types";
import type { CarbonFootprintData } from "../utils/carbonFootprint";
import { calculateCarbonFootprint } from "../utils/carbonFootprint";
import CarbonFootprint from "../components/CarbonFootprint";
import UpdateModal from "../components/UpdateModal";
import getCoordinatesFromAddress from "../utils/address_coordinates";
import '../index.css';
import OrderMap from '../components/OrderMap';
import {
  Container,
  Box,
  Typography,
  Card,
  CardContent,
  Button,
  Alert,
  CircularProgress,
  Skeleton,
  Divider,
} from '@mui/material';

export default function OrderPage() {
    const { id } = useParams<{ id: string }>(); 
    const [order, setOrder] = useState<OrderData | null>(null);
    const [orderHistory, setOrderHistory] = useState<OrderStatus[] | null>(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const [verifying, setVerifying] = useState(false);
    const [verificationResult, setVerificationResult] = useState<VerificationResult | null>(null);
    const [carbonFootprint, setCarbonFootprint] = useState<CarbonFootprintData | null>(null);
    const [showUpdateModal, setUpdateModal] = useState(false);
    const [address, setAddress] = useState("");
    const [statusMessage, setStatusMessage] = useState<string | null>(null);
    const [statusType, setStatusType] = useState<"success" | "error" | null>(null);
    const [isUpdating, setIsUpdating] = useState(false);


    useEffect(() => {
        const apiUrl = process.env.PUBLIC_API_URL || 'http://localhost:8080';
        
        fetch(`${apiUrl}/api/order/${id}`)
            .then(res => {
                if (!res.ok) setError("Could not load the order");
                return res.json();
            })
            .then(data => {
                const o = data.order as BackendOrder;
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
            
            fetch(`${apiUrl}/api/order/history/${id}`)
            .then(res => {
                if (!res.ok) setError("Could not load the order");
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
                history.sort((a, b) => b.timestamp.getTime() - a.timestamp.getTime());
                setOrderHistory(history);
            })
            .catch(() => {
                // History is optional, don't set error if it fails
                console.warn('Failed to load order history');
            });
    }, [id]);

    // Carbon footprint is calculated by OrderMap via callback
    const handleCarbonFootprintData = useCallback((data: { totalDistance: number; routeSegments: { distance: number; isAir: boolean }[] }) => {
        const footprint = calculateCarbonFootprint(data.totalDistance, data.routeSegments);
        setCarbonFootprint(footprint);
    }, []);

    const handleVerifyBlockchain = async () => {
        setVerifying(true);
        setVerificationResult(null);
        
        try {
            const apiUrl = process.env.PUBLIC_API_URL || 'http://localhost:8080';
            const response = await fetch(`${apiUrl}/api/order/verify/${id}`);
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

    const handleOrderUpdate = async () => {
        setIsUpdating(true);
        setStatusMessage(null); // reset previous message

        const result = await getCoordinatesFromAddress(address);

        if(!result){
            setStatusMessage("Invalid Address.");
            setStatusType("error");
            setIsUpdating(false);
            return
        }

        try {
            const apiUrl = process.env.PUBLIC_API_URL || "http://localhost:8080";
            const response = await fetch(`${apiUrl}/api/order/update`, {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify(
                {
                    order_id: Number(id),
                    delivery_address: address,
                    delivery_latitude: Number(result?.lat),
                    delivery_longitude: Number(result?.lon),
                }
            ),
            });

            if (!response.ok) {
                if(response.status == 403){
                    setStatusMessage("Update failed. Cannot change an order that is already shipped.");
                }
                else{
                    setStatusMessage("Update failed. Please try again.");
                }
                setStatusType("error");
                setIsUpdating(false);
            return;
            }

            const data = await response.json();

            if (data.error) {
                setStatusMessage("Update failed.");
                setStatusType("error");

            } else {
                setStatusMessage("Order updated successfully!");
                setStatusType("success");

                setOrder((prev) =>
                    prev
                    ? {
                        ...prev,
                        delivery_address: address,
                        delivery_latitude: Number(result?.lat),
                        delivery_longitude: Number(result?.lon),
                        }
                    : prev
                );

                // Close modal after successful update
                setTimeout(() => {
                    setUpdateModal(false);
                    setAddress("");
                }, 1500); // Delay to show success message
                
            }
        } catch (error) {
            setStatusMessage("Network error. Please try again.");
            setStatusType("error");
        }

        setIsUpdating(false);
    };




    if (loading) return (
        <Container maxWidth={false} sx={{ maxWidth: '64rem', py: 4 }}>
            {/* Map skeleton */}
            <Box sx={{ mb: 3, height: 384, bgcolor: 'action.hover', borderRadius: 2, display: 'flex', alignItems: 'center', justifyContent: 'center' }}>
                <CircularProgress />
            </Box>
            
            {/* Header skeleton */}
            <Card sx={{ mb: 3 }}>
                <CardContent>
                    <Box sx={{ display: 'flex', flexDirection: { xs: 'column', md: 'row' }, justifyContent: 'space-between', gap: 6, mb: 6 }}>
                        <Box sx={{ flex: 1, display: 'flex', flexDirection: 'column', gap: 4 }}>
                            <Skeleton variant="text" width="40%" height={40} />
                            <Box sx={{ display: 'flex', flexDirection: 'column', gap: 2 }}>
                                <Skeleton variant="text" width="80%" />
                                <Skeleton variant="text" width="70%" />
                                <Skeleton variant="text" width="90%" />
                            </Box>
                        </Box>
                        <Box sx={{ display: 'flex', flexDirection: 'column', gap: 4, minWidth: { md: '200px' } }}>
                            <Skeleton variant="rectangular" height={40} width="100%" />
                            <Box sx={{ display: 'flex', flexDirection: 'column', gap: 1, pt: 4, borderTop: '1px solid', borderColor: 'divider' }}>
                                <Skeleton variant="text" width="60%" height={24} />
                                <Skeleton variant="text" width="100%" height={40} />
                            </Box>
                        </Box>
                    </Box>
                </CardContent>
            </Card>
            
            {/* Products skeleton */}
            <Typography variant="h6" sx={{ mb: 2 }}>
                <Skeleton width="30%" />
            </Typography>
            <Box sx={{ space: 2, mb: 3 }}>
                {[1, 2].map(i => (
                    <Card key={i}>
                        <CardContent>
                            <Skeleton variant="text" width="50%" height={24} sx={{ mb: 1 }} />
                            <Skeleton variant="text" width="40%" />
                        </CardContent>
                    </Card>
                ))}
            </Box>
            
            {/* Tracking history skeleton */}
            <Typography variant="h6" sx={{ mb: 2 }}>
                <Skeleton width="30%" />
            </Typography>
            <Box>
                {[1, 2, 3].map(i => (
                    <Card key={i} sx={{ mb: 2 }}>
                        <CardContent>
                            <Skeleton variant="text" width="40%" height={24} sx={{ mb: 1 }} />
                            <Skeleton variant="text" width="80%" />
                            <Skeleton variant="text" width="60%" />
                        </CardContent>
                    </Card>
                ))}
            </Box>
        </Container>
    );
    
    if (error) return (
        <Container maxWidth="sm" sx={{ display: 'flex', alignItems: 'center', justifyContent: 'center', minHeight: '100vh' }}>
            <Box sx={{ textAlign: 'center' }}>
                <Typography variant="h5" sx={{ mb: 2, fontWeight: 600 }}>Unable to load order</Typography>
                <Typography variant="body2" color="textSecondary" sx={{ mb: 3 }}>{error}</Typography>
                <Button
                    variant="contained"
                    color="primary"
                    onClick={() => window.location.reload()}
                >
                    Try Again
                </Button>
            </Box>
        </Container>
    );
    
    if (!order) return (
        <Container maxWidth="sm" sx={{ display: 'flex', alignItems: 'center', justifyContent: 'center', minHeight: '100vh' }}>
            <Box sx={{ textAlign: 'center' }}>
                <Typography variant="h5" sx={{ mb: 2, fontWeight: 600 }}>Order not found</Typography>
                <Typography variant="body2" color="textSecondary">The order you&apos;re looking for doesn&apos;t exist or has been removed.</Typography>
            </Box>
        </Container>
    );


    return (
        <Container maxWidth={false} sx={{ maxWidth: '64rem', py: 4 }}>
            {/* Map */}
            {/* Wrapped lazy component in <Suspense> */}
            <Box sx={{ mb: 3 }}>
                <Suspense fallback={
                    <Box sx={{ width: '100%', height: 384, bgcolor: 'action.hover', borderRadius: 2, display: 'flex', alignItems: 'center', justifyContent: 'center' }}>
                        <CircularProgress />
                    </Box>
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
                            onCarbonFootprintData={handleCarbonFootprintData}
                        />
                    )}
                </Suspense>
            </Box>

            {/* Carbon Footprint */}
            {carbonFootprint && (
                <Box sx={{ mb: 3 }}>
                    <CarbonFootprint data={carbonFootprint} />
                </Box>
            )}

            {/* Header */}
            <Card sx={{ mb: 3 }}>
                <CardContent>
                    <Box sx={{ display: 'flex', flexDirection: { xs: 'column', md: 'row' }, justifyContent: 'space-between', gap: 6, mb: 6 }}>
                        <Box sx={{ flex: 1, display: 'flex', flexDirection: 'column', gap: 4 }}>
                            <Typography variant="h4" sx={{ fontWeight: 600 }}>Order {id}</Typography>
                            <Box sx={{ display: 'flex', flexDirection: 'column', gap: 2 }}>
                                <Box sx={{ display: 'flex', gap: 3 }}>
                                    <Typography variant="body2" color="textSecondary" sx={{ minWidth: 100 }}>Ordered</Typography>
                                    <Typography variant="body2">{new Date(order.created_at).toLocaleDateString()}</Typography>
                                </Box>
                                <Box sx={{ display: 'flex', gap: 3 }}>
                                    <Typography variant="body2" color="textSecondary" sx={{ minWidth: 100 }}>Est. Delivery</Typography>
                                    <Typography variant="body2">{new Date(order.delivery_estimate).toLocaleDateString()}</Typography>
                                </Box>
                                <Box sx={{ display: 'flex', gap: 3 }}>
                                    <Typography variant="body2" color="textSecondary" sx={{ minWidth: 100 }}>Delivery To</Typography>
                                    <Typography variant="body2">{order.delivery_address}</Typography>
                                </Box>
                            </Box>
                        </Box>
                        <Box sx={{ display: 'flex', flexDirection: 'column', gap: 4, minWidth: { md: '200px' } }}>
                            <Box sx={{ textAlign: { xs: 'left', md: 'right' } }}>
                                <Typography variant="caption" color="textSecondary" sx={{ display: 'block', mb: 1, textTransform: 'uppercase', letterSpacing: 1 }}>Tracking Code</Typography>
                                <Typography variant="body1" sx={{ fontFamily: 'monospace', fontWeight: 600 }}>{order.tracking_code}</Typography>
                            </Box>
                            <Box sx={{ textAlign: { xs: 'left', md: 'right' }, pt: 4, borderTop: '1px solid', borderColor: 'divider' }}>
                                <Typography variant="caption" color="textSecondary" sx={{ display: 'block', mb: 1, textTransform: 'uppercase', letterSpacing: 1 }}>Total</Typography>
                                <Typography variant="h5" sx={{ fontWeight: 600 }}>{order.price}€</Typography>
                            </Box>
                        </Box>
                    </Box>
                    <Box sx={{ display: 'flex', gap: 4, flexDirection: { xs: 'column', md: 'row' }, justifyContent: 'space-between', mt: 6 }}>
                        <Button
                            onClick={() => setUpdateModal(true)}
                            id="update"
                            disabled={showUpdateModal}
                            variant="contained"
                            color="primary"
                            size="small"
                            sx={{ width: { xs: '100%', md: 'auto' }, minWidth: { md: '200px' }, px: 4, py: 1, maxHeight: 40, mt: 1, fontSize: '0.875rem', fontWeight: 500, textTransform: 'none' }}
                        >
                            Update Address
                        </Button>
                        {/* Blockchain Verification */}
                        <Box sx={{ width: { xs: '100%', md: 'auto' }, minWidth: { md: '200px' }, maxWidth: { md: 80 } }}>
                            <Button
                                onClick={handleVerifyBlockchain}
                                id="verification"
                                disabled={verifying}
                                variant="contained"
                                color="primary"
                                size="small"
                                fullWidth
                                sx={{ px: 4, py: 1, maxHeight: 40, fontSize: '0.875rem', fontWeight: 500, textTransform: 'none', whiteSpace: 'nowrap' }}
                            >
                                {verifying ? (
                                    <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                                        <CircularProgress size={16} color="inherit" />
                                        Verifying...
                                    </Box>
                                ) : (
                                    'Verify on Blockchain'
                                )}
                            </Button>
                            
                            {verificationResult && (
                                <Box sx={{
                                    mt: 1.5,
                                    p: 1.5,
                                    borderRadius: 1,
                                    border: '1px solid',
                                    ...(verificationResult.status === 'VERIFIED' ? {
                                        bgcolor: 'rgba(76, 175, 80, 0.05)',
                                        borderColor: 'rgb(129, 199, 132)',
                                    } : verificationResult.status === 'ERROR' ? {
                                        bgcolor: 'rgba(244, 67, 54, 0.05)',
                                        borderColor: 'rgb(229, 119, 114)',
                                    } : {
                                        bgcolor: 'rgba(255, 152, 0, 0.05)',
                                        borderColor: 'rgb(255, 167, 38)',
                                    })
                                }}>
                                    <Box sx={{ display: 'flex', alignItems: 'center', gap: 1, mb: 0.75 }}>
                                        <Typography variant="subtitle2" sx={{
                                            fontWeight: 600,
                                            ...(verificationResult.status === 'VERIFIED' ? {
                                                color: '#2e7d32',
                                            } : verificationResult.status === 'ERROR' ? {
                                                color: '#c62828',
                                            } : {
                                                color: '#e65100',
                                            })
                                        }}>
                                            {verificationResult.status.replace('_', ' ')}
                                        </Typography>
                                    </Box>
                                    <Typography variant="caption" sx={{
                                        display: 'block',
                                        mb: 0.75,
                                        ...(verificationResult.status === 'VERIFIED' ? {
                                            color: '#558b2f',
                                        } : verificationResult.status === 'ERROR' ? {
                                            color: '#b71c1c',
                                        } : {
                                            color: '#bf360c',
                                        })
                                    }}>
                                        {verificationResult.message}
                                    </Typography>
                                    {verificationResult.status !== 'ERROR' && (
                                        <>
                                            <Typography variant="caption" sx={{
                                                display: 'block',
                                                mb: 0.75,
                                                mt: 1,
                                                color: 'text.secondary'
                                            }}>
                                                {verificationResult.verified_updates}/{verificationResult.total_updates} updates verified
                                            </Typography>
                                            {verificationResult.contract_address && (
                                                <Typography variant="caption" sx={{
                                                    display: 'block',
                                                    mb: 0.75,
                                                    mt: 0.5,
                                                    wordBreak: 'break-word',
                                                    color: 'text.secondary'
                                                }}>
                                                    <strong>Contract:</strong> {verificationResult.contract_address}
                                                </Typography>
                                            )}
                                            {verificationResult.transaction_hashes && verificationResult.transaction_hashes.length > 0 && (
                                                <Typography variant="caption" sx={{
                                                    display: 'block',
                                                    mt: 1,
                                                    color: 'text.secondary'
                                                }}>
                                                    ✓ Blockchain links shown in tracking history below
                                                </Typography>
                                            )}
                                        </>
                                    )}
                                </Box>
                            )}
                        </Box>
                    </Box>
                </CardContent>
            </Card>

            {/*Modal for to update the order information*/ }
            <UpdateModal show={showUpdateModal} onClose={() => setUpdateModal(false)} isUpdating={isUpdating} onUpdate={handleOrderUpdate}>
                <Typography variant="h6" sx={{ mb: 2, fontWeight: 600 }}>Update Delivery Address</Typography>
                   {/* Feedback message injected as children */}
                    {statusMessage && (
                        <Alert severity={statusType === "success" ? "success" : "error"} sx={{ mb: 2 }}>
                            {statusMessage}
                        </Alert>
                    )}
                    {/* Text box for delivery address */}
                    <input
                    type="text"
                    value={address}
                    id = "delivery-address"
                    onChange={(e) => setAddress(e.target.value)}
                    placeholder="Enter new delivery address"
                    className="w-full px-3 py-2 border border-gray-300 rounded-lg mb-4 focus:outline-none focus:ring-2 focus:ring-blue-500 text-gray-900 dark:text-white"
                    />
            </UpdateModal>

            {/* Products */}
            <Box sx={{ mb: 6 }}>
                <Typography variant="h5" sx={{ mb: 2, fontWeight: 600 }}>Order Items</Typography>
                {order.products && order.products.length > 0 ? (
                    <Box sx={{ display: 'flex', flexDirection: 'column', gap: 2 }}>
                        {order.products.map((p, idx) => (
                            <Card key={p.product_id || idx}>
                                <CardContent>
                                    <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start', gap: 2, mb: 2 }}>
                                        <Box sx={{ flex: 1 }}>
                                            <Typography variant="subtitle1" sx={{ fontWeight: 600, mb: 1 }}>{p.name}</Typography>
                                            <Typography variant="caption" color="textSecondary">Product ID: {p.product_id}</Typography>
                                        </Box>
                                        <Box sx={{ textAlign: 'right' }}>
                                            <Box sx={{ display: 'flex', alignItems: 'baseline', gap: 1, justifyContent: 'flex-end' }}>
                                                <Typography variant="h6" sx={{ fontWeight: 600, fontSize: '1.25rem' }}>{p.price}€</Typography>
                                            </Box>
                                            <Typography variant="caption" color="textSecondary">per unit</Typography>
                                        </Box>
                                    </Box>
                                    <Divider sx={{ my: 2 }} />
                                    <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                                        <Box sx={{ display: 'flex', gap: 1 }}>
                                            <Typography variant="caption" color="textSecondary">Quantity:</Typography>
                                            <Typography variant="caption" sx={{ fontWeight: 600 }}>{p.quantity}</Typography>
                                        </Box>
                                        <Box sx={{ textAlign: 'right', display: 'flex', alignItems: 'center', gap: 1, justifyContent: 'flex-end' }}>
                                            <Typography variant="caption" color="textSecondary">Subtotal:</Typography>
                                            <Typography variant="h6" sx={{ fontWeight: 600, fontSize: '1.125rem' }}>{(p.price * p.quantity).toFixed(2)}€</Typography>
                                        </Box>
                                    </Box>
                                </CardContent>
                            </Card>
                        ))}
                        
                        {/* Total summary */}
                        <Card sx={{ bgcolor: 'action.hover' }}>
                            <CardContent>
                                <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                                    <Box>
                                        <Typography variant="body2" sx={{ fontWeight: 600 }}>Total</Typography>
                                        <Typography variant="caption" color="textSecondary">{order.products.length} item{order.products.length !== 1 ? 's' : ''}</Typography>
                                    </Box>
                                    <Typography variant="h5" sx={{ fontWeight: 600 }}>
                                        {order.price}€
                                    </Typography>
                                </Box>
                            </CardContent>
                        </Card>
                    </Box>
                ) : (
                    <Alert severity="info">No products in this order.</Alert>
                )}
            </Box>

            {/* Tracking Steps */}
            <Box>
                <Typography variant="h5" sx={{ mb: 2, fontWeight: 600 }}>Tracking History</Typography>
                {orderHistory && orderHistory.length > 0 ? (
                    <Box sx={{ display: 'flex', flexDirection: 'column', gap: 0 }}>
                        {orderHistory.map((s, idx) => {
                            const isFirst = idx === 0;
                            const isLast = idx === orderHistory.length - 1;
                            
                            // Define status colors using theme
                            const statusInfo: Record<string, { color: string; severity: 'info' | 'success' | 'warning' | 'error' }> = {
                                'PROCESSING': { color: '#1976d2', severity: 'info' },
                                'SHIPPED': { color: '#1976d2', severity: 'info' },
                                'IN TRANSIT': { color: '#1976d2', severity: 'info' },
                                'OUT FOR DELIVERY': { color: '#1976d2', severity: 'info' },
                                'DELIVERED': { color: '#388e3c', severity: 'success' },
                                'CANCELLED': { color: '#d32f2f', severity: 'error' },
                                'RETURNED': { color: '#f57c00', severity: 'warning' },
                                'FAILED DELIVERY': { color: '#d32f2f', severity: 'error' }
                            };
                            
                            const status = statusInfo[s.order_status] || { color: '#9e9e9e', severity: 'info' as const };
                            
                            // Get the blockchain transaction hash for this status update
                            const txHash = verificationResult?.transaction_hashes?.[idx];
                            
                            return (
                                <Box key={idx} sx={{ display: 'flex', mb: isLast ? 0 : 0, alignItems: 'stretch' }}>
                                    {/* Left timeline column */}
                                    <Box sx={{ display: 'flex', flexDirection: 'column', alignItems: 'center', width: 48, flexShrink: 0 }}>
                                        {/* Top connecting line - only if not first */}
                                        {!isFirst && (
                                            <Box sx={{ width: 2, height: 24, bgcolor: 'divider' }} />
                                        )}
                                        
                                        {/* The dot - larger for the first (current state) one */}
                                        <Box
                                            sx={{
                                                width: isFirst ? 32 : 24,
                                                height: isFirst ? 32 : 24,
                                                borderRadius: '50%',
                                                bgcolor: status.color,
                                                flexShrink: 0,
                                                border: 4,
                                                borderColor: 'background.paper',
                                            }}
                                        />
                                        
                                        {/* Bottom connecting line - only if not last */}
                                        {!isLast && (
                                            <Box sx={{ width: 2, flex: 1, minHeight: 24, bgcolor: 'divider' }} />
                                        )}
                                    </Box>
                                    
                                    {/* Content card */}
                                    <Box sx={{ flex: 1, pb: isLast ? 0 : 2, pl: 2 }}>
                                        <Card>
                                            <CardContent>
                                                <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start', mb: 1 }}>
                                                    <Typography variant="subtitle1" sx={{ fontWeight: 600 }}>{s.order_status}</Typography>
                                                    <Typography variant="caption" color="textSecondary" sx={{ whiteSpace: 'nowrap', ml: 2 }}>
                                                        {new Date(s.timestamp).toLocaleString()}
                                                    </Typography>
                                                </Box>
                                                
                                                <Typography variant="body2" color="textSecondary" sx={{ mb: 1 }}>{s.note}</Typography>
                                                
                                                <Typography variant="body2" color="textSecondary">
                                                    {s.order_location}
                                                </Typography>
                                                
                                                {s.Storage && (
                                                    <>
                                                        <Divider sx={{ my: 1.5 }} />
                                                        <Typography variant="subtitle2" sx={{ fontWeight: 600, mb: 0.5 }}>{s.Storage.Name}</Typography>
                                                        <Typography variant="caption" color="textSecondary">{s.Storage.Address}</Typography>
                                                    </>
                                                )}
                                                
                                                {/* Blockchain verification link */}
                                                {txHash && (
                                                    <>
                                                        <Divider sx={{ my: 1.5 }} />
                                                        <Typography variant="caption" sx={{ display: 'block', color: 'success.main', fontWeight: 600, mb: 0.5 }}>✓ Verified on Blockchain</Typography>
                                                        <Typography
                                                            component="a"
                                                            href={`https://sepolia.etherscan.io/tx/${txHash}`}
                                                            target="_blank"
                                                            rel="noopener noreferrer"
                                                            id={`link-${idx}`}
                                                            variant="caption"
                                                            sx={{
                                                                display: 'block',
                                                                color: 'primary.main',
                                                                textDecoration: 'none',
                                                                '&:hover': { textDecoration: 'underline' },
                                                                fontFamily: 'monospace',
                                                                wordBreak: 'break-all',
                                                            }}
                                                            title="View transaction on Etherscan"
                                                        >
                                                            {txHash}
                                                        </Typography>
                                                    </>
                                                )}
                                            </CardContent>
                                        </Card>
                                    </Box>
                                </Box>
                            );
                        })}
                    </Box>
                ) : (
                    <Alert severity="info">No tracking history available yet.</Alert>
                )}
            </Box>
        </Container>
    );
}