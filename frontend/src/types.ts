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
    id?: number;
    tracking_code: string;
    delivery_estimate: string;
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

// Backend API response types

interface BackendOrderProduct {
    Product_ID: number;
    Quantity: number;
    Product_Name_At_Purchase: string;
    Product_Price_At_Purchase: number;
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
    Id?: number;
    Tracking_Code: string;
    Delivery_Estimate: string;
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

interface VerificationResult {
    verified: boolean;
    total_updates: number;
    verified_updates: number;
    blockchain_hashes: number;
    status: string;
    message: string;
    mismatches?: string[];
    transaction_hashes?: string[];
    contract_address?: string;
}

export type {OrderData, OrderProduct, OrderStatus, Storage, VerificationResult, BackendOrderStatus, BackendOrder, BackendStorage, BackendOrderProduct}