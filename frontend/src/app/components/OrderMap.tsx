"use client";

import { MapContainer, TileLayer, Marker, Popup, Polyline } from 'react-leaflet';
import { OrderStatus } from '@/app/types';
import 'leaflet/dist/leaflet.css';
import L from 'leaflet';
import { useState, useEffect } from 'react';
import RoutingMachine from './RoutingMachine';

// Fix for default marker icons in Next.js
delete (L.Icon.Default.prototype as any)._getIconUrl;
L.Icon.Default.mergeOptions({
    iconRetinaUrl: 'https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.7.1/images/marker-icon-2x.png',
    iconUrl: 'https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.7.1/images/marker-icon.png',
    shadowUrl: 'https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.7.1/images/marker-shadow.png',
});

// Custom icons
const deliveryIcon = L.icon({
    iconUrl: 'https://raw.githubusercontent.com/pointhi/leaflet-color-markers/master/img/marker-icon-2x-green.png',
    shadowUrl: 'https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.7.1/images/marker-shadow.png',
    iconSize: [25, 41],
    iconAnchor: [12, 41],
    popupAnchor: [1, -34],
    shadowSize: [41, 41]
});
const sellerIcon = L.icon({
    iconUrl: 'https://raw.githubusercontent.com/pointhi/leaflet-color-markers/master/img/marker-icon-2x-red.png',
    shadowUrl: 'https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.7.1/images/marker-shadow.png',
    iconSize: [25, 41],
    iconAnchor: [12, 41],
    popupAnchor: [1, -34],
    shadowSize: [41, 41]
});

interface OrderMapProps {
    orderHistory: OrderStatus[];
    deliveryAddress?: string;
    deliveryLatitude?: number;
    deliveryLongitude?: number;
    sellerAddress?: string;
    sellerLatitude?: number;
    sellerLongitude?: number;
}

export default function OrderMap({
    orderHistory,
    deliveryAddress,
    deliveryLatitude,
    deliveryLongitude,
    sellerAddress,
    sellerLatitude,
    sellerLongitude
}: OrderMapProps) {
    const [mapInstance, setMapInstance] = useState<L.Map | null>(null);

    const locationsWithCoords = orderHistory.filter(h => h.Storage);
    const routeCoordinates: [number, number][] = locationsWithCoords.map(h => [
        h.Storage!.Latitude,
        h.Storage!.Longitude
    ]);
    const sellerCoords: [number, number] | null =
        sellerLatitude && sellerLongitude ? [sellerLatitude, sellerLongitude] : null;
    const deliveryCoords: [number, number] | null =
        deliveryLatitude && deliveryLongitude ? [deliveryLatitude, deliveryLongitude] : null;

    // Center map on all points
    const allCoords: [number, number][] = [
        ...(sellerCoords ? [sellerCoords] : []),
        ...routeCoordinates,
        ...(deliveryCoords ? [deliveryCoords] : [])
    ];
    const center: [number, number] = allCoords.length > 0
        ? [
            allCoords.reduce((sum, coord) => sum + coord[0], 0) / allCoords.length,
            allCoords.reduce((sum, coord) => sum + coord[1], 0) / allCoords.length
        ]
        : [39.5, -8.0];


    const latestStatus = orderHistory.length > 0 ? orderHistory[orderHistory.length - 1] : null;
    const currentStatus = latestStatus?.order_status;
    
    let routeColor: 'green' | 'red' | 'blue' = 'blue';
    let shouldZoomOut = false;

    if (currentStatus === 'DELIVERED') {
        routeColor = 'green';
        shouldZoomOut = true;
    } else if (
        currentStatus === 'CANCELLED' || 
        currentStatus === 'RETURNED' || 
        currentStatus === 'FAILED DELIVERY'
    ) {
        routeColor = 'red';
        shouldZoomOut = true;
    }

    const partialRoute: [number, number][] = [];

    // 1. Add seller
    if (sellerCoords) {
        partialRoute.push(sellerCoords);
    }

    // 2. Add all *visited* storage locations from history
    partialRoute.push(...routeCoordinates);

    // 3. ONLY add delivery location if the order is marked as delivered
    if (currentStatus === 'DELIVERED' && deliveryCoords) {
        partialRoute.push(deliveryCoords);
    }

    useEffect(() => {
        if (!mapInstance) return;

        if (shouldZoomOut) {
            if (allCoords.length > 0) {
                mapInstance.fitBounds(allCoords, { padding: [50, 50] });
            }
        } else {
            if (partialRoute.length > 0) {
                const lastKnownLocation = partialRoute[partialRoute.length - 1];
                mapInstance.flyTo(lastKnownLocation, 12, { duration: 1.5 });
            }
        }
    }, [mapInstance, shouldZoomOut, partialRoute, allCoords]);

    const routeSegments: [[number, number], [number, number]][] = [];
    for (let i = 0; i < partialRoute.length - 1; i++) {
        routeSegments.push([partialRoute[i], partialRoute[i + 1]]);
    }

    // 500km threshold to consider a route as sea route
    const SEA_ROUTE_THRESHOLD_METERS = 500 * 1000;

    if (locationsWithCoords.length === 0 && !sellerCoords) {
        return (
            <div className="w-full h-96 bg-gray-200 rounded-lg flex items-center justify-center">
                <p className="text-gray-500">No location data available</p>
            </div>
        );
    }

    const handleZoomToDelivery = () => {
        if (deliveryCoords && mapInstance) {
            mapInstance.flyTo(deliveryCoords, 12, { duration: 1.5 });
        }
    };
    const handleZoomToSeller = () => {
        if (sellerCoords && mapInstance) {
            mapInstance.flyTo(sellerCoords, 12, { duration: 1.5 });
        }
    };

    return (
        <MapContainer
            center={center}
            zoom={7}
            style={{ height: '500px', width: '100%', borderRadius: '0.5rem' }}
            ref={setMapInstance}
        >
            <TileLayer
                attribution='&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
                url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
            />

            {routeSegments.map((segment, idx) => {
                // Calculate straight-line distance (in meters)
                const distance = L.latLng(segment[0]).distanceTo(L.latLng(segment[1]));

                if (distance > SEA_ROUTE_THRESHOLD_METERS) {
                    // It's a sea route, draw a straight line
                    return (
                        <Polyline
                            key={idx}
                            positions={segment}
                            color={routeColor}
                            weight={5}
                            opacity={0.7}
                            dashArray="5, 10"
                        />
                    );
                } else {
                    // It's a land route, use the routing machine
                    return (
                        <RoutingMachine
                            key={idx}
                            waypoints={segment}
                            routeColor={routeColor}
                        />
                    );
                }
            })}

            {/* Seller marker */}
            {sellerCoords && (
                <Marker position={sellerCoords} icon={sellerIcon}>
                    <Popup>
                        <div className="text-sm">
                            <p className="font-bold">Seller</p>
                            <p className="text-xs mt-1">{sellerAddress}</p>
                        </div>
                    </Popup>
                </Marker>
            )}

            {/* Storage markers */}
            {locationsWithCoords.map((history, idx) => (
                <Marker
                    key={idx}
                    position={[history.Storage!.Latitude, history.Storage!.Longitude]}
                >
                    <Popup>
                        <div className="text-sm">
                            <p className="font-bold">{history.Storage!.Name}</p>
                            <p>{history.Storage!.Address}</p>
                            <p className="mt-2 text-xs text-gray-600">
                                Status: {history.order_status}
                            </p>
                            <p className="text-xs text-gray-600">{history.note}</p>
                            <p className="text-xs text-gray-500 mt-1">
                                {new Date(history.timestamp).toLocaleString()}
                            </p>
                        </div>
                    </Popup>
                </Marker>
            ))}

            {/* Delivery marker */}
            {deliveryCoords && (
                <Marker position={deliveryCoords} icon={deliveryIcon}>
                    <Popup>
                        <div className="text-sm">
                            <p className="font-bold text-green-600">Delivery Destination</p>
                            <p className="text-xs mt-1">{deliveryAddress}</p>
                        </div>
                    </Popup>
                </Marker>
            )}

            {/* Info overlays */}
            {sellerAddress && (
                <div className="leaflet-bottom leaflet-left" style={{ pointerEvents: 'none' }}>
                    <div
                        className="bg-white p-3 m-4 rounded-lg shadow-lg cursor-pointer hover:shadow-xl transition-shadow"
                        style={{ pointerEvents: 'auto' }}
                        onClick={handleZoomToSeller}
                    >
                        <p className="text-xs font-semibold text-gray-700">Seller:</p>
                        <p className="text-xs text-gray-600">{sellerAddress}</p>
                        {sellerCoords && (
                            <p className="text-xs text-red-600 mt-1 font-semibold">
                                📍 Click to zoom to seller
                            </p>
                        )}
                    </div>
                </div>
            )}
            {deliveryAddress && (
                <div className="leaflet-bottom leaflet-right" style={{ pointerEvents: 'none' }}>
                    <div
                        className="bg-white p-3 m-4 rounded-lg shadow-lg cursor-pointer hover:shadow-xl transition-shadow"
                        style={{ pointerEvents: 'auto' }}
                        onClick={handleZoomToDelivery}
                    >
                        <p className="text-xs font-semibold text-gray-700">Delivery To:</p>
                        <p className="text-xs text-gray-600">{deliveryAddress}</p>
                        {deliveryCoords && (
                            <p className="text-xs text-green-600 mt-1 font-semibold">
                                📍 Click to zoom to destination
                            </p>
                        )}
                    </div>
                </div>
            )}
        </MapContainer>
    );
}