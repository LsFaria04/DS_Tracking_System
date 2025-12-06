import { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';
import type { OrderData, BackendOrder, BackendOrderProduct, BackendOrderStatus } from '../types';
import '../index.css';
import {
  Container,
  Box,
  Typography,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  Card,
  CardContent,
  Skeleton,
  Alert,
  useTheme,
} from '@mui/material';

export default function OrdersPage() {
    const [orders, setOrders] = useState<OrderData[] | null>(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const [order_by, setOrderBy] = useState<string>("newest");
    const [statusFilter, setStatusFilter] = useState<string>("all");
    const theme = useTheme();

    const apiUrl = process.env.PUBLIC_API_URL || 'http://localhost:8080';

    useEffect(() => {
        setLoading(true);
        handleOrders();
    }, [order_by, statusFilter]);

    function handleOrders(){
        fetch(`${apiUrl}/api/orders?order_by=${order_by}`)
            .then((res) => {
                if (!res.ok) setError("Could not load the orders");
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
                    statusHistory:ord.Updates?.map((p: BackendOrderStatus) => ({
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
        <Container maxWidth={false} sx={{ maxWidth: '64rem', py: 4 }}>
            <Typography variant="h4" sx={{ mb: 4, fontWeight: 600 }}>Orders</Typography>
            {/* Status filter skeleton */}
            {/* Sort selector skeleton */}
            <Box sx={{ display: 'flex', gap: 2, mb: 4, flexDirection: { xs: 'column', md: 'row' } }}>
                <Skeleton variant="rectangular" width="100%" height={40} sx={{ maxWidth: 200 }} />
                <Skeleton variant="rectangular" width="100%" height={40} sx={{ maxWidth: 200 }} />
            </Box>
            <Box sx={{ display: 'flex', flexDirection: 'column', gap: 2 }}>
                {[1, 2, 3].map((i) => (
                    <Card key={i}>
                        <CardContent>
                            <Skeleton variant="text" width="30%" height={32} sx={{ mb: 1 }} />
                            <Skeleton variant="text" width="60%" height={20} />
                        </CardContent>
                    </Card>
                ))}
            </Box>
        </Container>
    );

    if (error) return (
        <Container maxWidth={false} sx={{ maxWidth: '64rem', display: 'flex', alignItems: 'center', justifyContent: 'center', minHeight: '100vh' }}>
            <Box sx={{ textAlign: 'center' }}>
                <Typography variant="h5" sx={{ mb: 2, fontWeight: 600 }}>Unable to load orders</Typography>
                <Typography variant="body2" color="textSecondary" sx={{ mb: 3 }}>{error}</Typography>
            </Box>
        </Container>
    );

    return (
        <Container maxWidth={false} sx={{ maxWidth: '64rem', py: 4 }}>
            <Typography variant="h4" sx={{ mb: 4, fontWeight: 600 }}>Orders</Typography>
            {/* Status filter */}
            {/* Sort selector */}
            <Box sx={{ display: 'flex', gap: 2, mb: 4, flexDirection: { xs: 'column', md: 'row' } }}>
                <FormControl sx={{ minWidth: 200 }}>
                    <InputLabel>Status</InputLabel>
                    <Select
                        value={statusFilter}
                        label="Status"
                        onChange={(e) => setStatusFilter(e.target.value)}
                    >
                        <MenuItem value="all">All statuses</MenuItem>
                        <MenuItem value="PROCESSING">Processing</MenuItem>
                        <MenuItem value="SHIPPED">Shipped</MenuItem>
                        <MenuItem value="IN TRANSIT">In Transit</MenuItem>
                        <MenuItem value="OUT FOR DELIVERY">Out for Delivery</MenuItem>
                        <MenuItem value="DELIVERED">Delivered</MenuItem>
                        <MenuItem value="CANCELLED">Cancelled</MenuItem>
                    </Select>
                </FormControl>

                <FormControl sx={{ minWidth: 200 }}>
                    <InputLabel>Sort by</InputLabel>
                    <Select
                        value={order_by}
                        label="Sort by"
                        onChange={(e) => setOrderBy(e.target.value)}
                    >
                        <MenuItem value="newest">Newest</MenuItem>
                        <MenuItem value="oldest">Oldest</MenuItem>
                    </Select>
                </FormControl>
            </Box>

            {orders && orders.length > 0 ? (
                <Box sx={{ display: 'flex', flexDirection: 'column', gap: 2 }}>
                    {orders.map((o) => (
                        <Card
                            key={o.id}
                            id={`order-${o.id ? o.id.toString() : "1"}`}
                            component={Link}
                            to={`/order/${o.id}`}
                            sx={{
                                textDecoration: 'none',
                                border: `1px solid ${theme.palette.divider}`,
                                transition: 'all 0.2s',
                                '&:hover': {
                                    borderColor: theme.palette.mode === 'dark' ? '#484848' : '#d0d0d0',
                                    bgcolor: theme.palette.mode === 'dark' ? '#2a2a2a' : '#f5f5f5',
                                    boxShadow: 2,
                                },
                            }}
                        >
                            <CardContent>
                                <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start' }}>
                                    <Box>
                                        <Typography variant="subtitle1" sx={{ fontWeight: 600 }}>{o.tracking_code}</Typography>
                                        <Typography variant="body2" color="textSecondary">{o.delivery_address}</Typography>
                                    </Box>
                                    <Box sx={{ textAlign: 'right' }}>
                                        <Typography variant="subtitle1" sx={{ fontWeight: 600 }}>{o.price}€</Typography>
                                        <Typography variant="caption" color="textSecondary">{o.products?.length ?? 0} item{o.products?.length !== 1 ? 's' : ''}</Typography>
                                    </Box>
                                </Box>
                            </CardContent>
                        </Card>
                    ))}
                </Box>
            ) : (
                <Alert severity="info">No orders found.</Alert>
            )}
        </Container>
    );
}