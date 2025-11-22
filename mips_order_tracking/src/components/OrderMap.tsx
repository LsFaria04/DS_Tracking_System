import { MapContainer, TileLayer, Marker, Popup, Polyline } from 'react-leaflet';
import type { OrderStatus } from '../types';
import 'leaflet/dist/leaflet.css';
import * as L from 'leaflet';
import { useState, useEffect, useMemo, useCallback } from 'react';
import RoutingMachine from './RoutingMachine';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faTruck } from '@fortawesome/free-solid-svg-icons';
import { getTomTomTrafficData } from '../utils/trafficApi';

// Fix for default marker icons
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

    const [totalDuration, setTotalDuration] = useState<number>(0);
    const [totalDistance, setTotalDistance] = useState<number>(0);
    const [routesCalculated, setRoutesCalculated] = useState<Set<number>>(new Set());
    const [estimatedTrafficDelay, setEstimatedTrafficDelay] = useState<number>(0);

    // Sort history by timestamp ascending (old -> new) so the route is built chronologically
    const sortedByTimeAsc = useMemo(() => 
        [...orderHistory].sort((a, b) => a.timestamp.getTime() - b.timestamp.getTime()),
        [orderHistory]
    );

    // Storages visited in chronological order
    const locationsWithCoords = useMemo(() => 
        sortedByTimeAsc.filter(h => h.Storage),
        [sortedByTimeAsc]
    );
    
    const routeCoordinates: [number, number][] = useMemo(() => 
        locationsWithCoords.map(h => [h.Storage!.Latitude, h.Storage!.Longitude]),
        [locationsWithCoords]
    );

    // Only accept numeric lat/lon (0 is valid) — check with typeof
    const sellerCoords: [number, number] | null = useMemo(() =>
        (typeof sellerLatitude === 'number' && typeof sellerLongitude === 'number') 
            ? [sellerLatitude, sellerLongitude] : null,
        [sellerLatitude, sellerLongitude]
    );
    
    const deliveryCoords: [number, number] | null = useMemo(() =>
        (typeof deliveryLatitude === 'number' && typeof deliveryLongitude === 'number') 
            ? [deliveryLatitude, deliveryLongitude] : null,
        [deliveryLatitude, deliveryLongitude]
    );

    // Center map on all points
    const allCoords: [number, number][] = useMemo(() => [
        ...(sellerCoords ? [sellerCoords] : []),
        ...routeCoordinates,
        ...(deliveryCoords ? [deliveryCoords] : [])
    ], [sellerCoords, routeCoordinates, deliveryCoords]);

    const center: [number, number] = useMemo(() => 
        allCoords.length > 0
            ? [
                allCoords.reduce((sum, coord) => sum + coord[0], 0) / allCoords.length,
                allCoords.reduce((sum, coord) => sum + coord[1], 0) / allCoords.length
            ]
            : [39.5, -8.0],
        [allCoords]
    );


    // Determine latest status using the sorted history (most recent timestamp)
    const latestStatus = sortedByTimeAsc.length > 0 ? sortedByTimeAsc[sortedByTimeAsc.length - 1] : null;
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

    const partialRoute: [number, number][] = useMemo(() => {
        const route: [number, number][] = [];
        
        // 1. Add seller
        if (sellerCoords) {
            route.push(sellerCoords);
        }

        // 2. Add all *visited* storage locations from history
        route.push(...routeCoordinates);

        // 3. ONLY add delivery location if the order is marked as delivered
        if (currentStatus === 'DELIVERED' && deliveryCoords) {
            route.push(deliveryCoords);
        }

        return route;
    }, [sellerCoords, routeCoordinates, currentStatus, deliveryCoords]);

    const routeSegments: [[number, number], [number, number]][] = useMemo(() => {
        const segments: [[number, number], [number, number]][] = [];
        for (let i = 0; i < partialRoute.length - 1; i++) {
            segments.push([partialRoute[i], partialRoute[i + 1]]);
        }
        return segments;
    }, [partialRoute]);

    // Reset metrics when route changes
    useEffect(() => {
        setTotalDistance(0);
        setTotalDuration(0);
        setRoutesCalculated(new Set());
        setEstimatedTrafficDelay(0);
    }, [routeSegments]);

    // Callback to accumulate route data - only once per segment
    const handleRouteFound = useCallback((segmentIndex: number) => 
        (summary: { distance: number; duration: number }) => {
            setRoutesCalculated(prev => {
                // Only add if we haven't calculated this segment yet
                if (prev.has(segmentIndex)) return prev;
                
                const newSet = new Set(prev);
                newSet.add(segmentIndex);
                
                setTotalDistance(prevDist => prevDist + summary.distance);
                setTotalDuration(prevDur => prevDur + summary.duration);
                
                return newSet;
            });
        }, 
    []);

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

    // 500km threshold to consider a route as sea route
    const SEA_ROUTE_THRESHOLD_METERS = 500 * 1000;

    const landSegmentCount = useMemo(() => 
        routeSegments.filter(seg => {
            const distance = L.latLng(seg[0]).distanceTo(L.latLng(seg[1]));
            return distance <= SEA_ROUTE_THRESHOLD_METERS;
        }).length,
        [routeSegments]
    );

    // Calculate traffic delay when all routes are calculated
    useEffect(() => {
        // Only calculate for orders in active transit
        if (currentStatus === 'DELIVERED' || 
            currentStatus === 'CANCELLED' || 
            currentStatus === 'RETURNED' || 
            currentStatus === 'FAILED DELIVERY' || 
            currentStatus === 'PROCESSING') {
            return;
        }
        
        if (routesCalculated.size === landSegmentCount && totalDuration > 0 && totalDistance > 0) {
            const tomtomKey = process.env.PUBLIC_TOMTOM_API_KEY;
            
            if (tomtomKey && partialRoute.length >= 2 && deliveryCoords) {
                const currentLocation = partialRoute[partialRoute.length - 1];
                
                // Check if current location is different from delivery destination
                if (currentLocation[0] !== deliveryCoords[0] || currentLocation[1] !== deliveryCoords[1]) {
                    getTomTomTrafficData(currentLocation, deliveryCoords, tomtomKey)
                        .then(data => {
                            if (data && data.trafficDelay > 0) {
                                setEstimatedTrafficDelay(data.trafficDelay);
                            } else {
                                setEstimatedTrafficDelay(0);
                            }
                        })
                        .catch(error => {
                            console.warn('Failed to fetch traffic data:', error);
                            setEstimatedTrafficDelay(0);
                        });
                } else {
                    setEstimatedTrafficDelay(0);
                }
            } else {
                setEstimatedTrafficDelay(0);
            }
        }
    }, [routesCalculated.size, landSegmentCount, totalDuration, totalDistance, partialRoute, deliveryCoords, currentStatus]);

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
        <div className="relative">
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
                    const isSeaRoute = distance > SEA_ROUTE_THRESHOLD_METERS;

                    if (isSeaRoute) {
                        // It's a sea route, draw a straight line
                        return (
                            <Polyline
                                key={`sea-${idx}`}
                                positions={segment}
                                pathOptions={{
                                    color: routeColor,
                                    weight: 5,
                                    opacity: 0.7,
                                    dashArray: '5, 10'
                                }}
                            />
                        );
                    } else {
                        // It's a land route, use the routing machine
                        return (
                            <RoutingMachine
                                key={`land-${idx}`}
                                waypoints={segment}
                                routeColor={routeColor}
                                onRouteFound={handleRouteFound(idx)}
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

            {/* Show estimate only when all land routes are calculated and order is in active transit */}
            {totalDuration > 0 && routesCalculated.size === landSegmentCount && 
                currentStatus !== 'DELIVERED' && 
                currentStatus !== 'CANCELLED' && 
                currentStatus !== 'RETURNED' && 
                currentStatus !== 'FAILED DELIVERY' && 
                currentStatus !== 'PROCESSING' && (
                <div className="absolute top-4 right-4 bg-white dark:bg-gray-900 border border-gray-200 dark:border-gray-800 rounded-lg shadow-lg p-4 max-w-xs z-[1000]">
                    <div className="flex items-center gap-2 mb-3">
                        <FontAwesomeIcon icon={faTruck} className="text-xl text-blue-600 dark:text-blue-400" />
                        <h3 className="font-semibold text-gray-900 dark:text-white">Estimated Delivery</h3>
                    </div>
                    
                    <div className="space-y-3 text-sm">
                        <div className="flex justify-between items-center">
                            <span className="text-gray-600 dark:text-gray-400">Travel time:</span>
                            <span className="font-medium text-gray-900 dark:text-white text-lg">
                                {(() => {
                                    const totalMinutes = Math.round(totalDuration / 60);
                                    const days = Math.floor(totalMinutes / (24 * 60));
                                    const hours = Math.floor((totalMinutes % (24 * 60)) / 60);
                                    const minutes = totalMinutes % 60;
                                    
                                    if (days > 0) {
                                        return `~${days}d ${hours}h`;
                                    } else if (hours > 0) {
                                        return `~${hours}h ${minutes}min`;
                                    } else {
                                        return `~${minutes}min`;
                                    }
                                })()}
                            </span>
                        </div>
                        
                        {estimatedTrafficDelay > 0 && (
                            <div className="flex justify-between items-center">
                                <span className="text-orange-600 dark:text-orange-400 text-xs">+ Traffic delay:</span>
                                <span className="font-medium text-orange-600 dark:text-orange-400">
                                    {(() => {
                                        const delayMinutes = Math.round(estimatedTrafficDelay / 60);
                                        const hours = Math.floor(delayMinutes / 60);
                                        const minutes = delayMinutes % 60;
                                        
                                        if (hours > 0) {
                                            return `~${hours}h ${minutes}min`;
                                        } else {
                                            return `~${minutes}min`;
                                        }
                                    })()}
                                </span>
                            </div>
                        )}
                        
                        {estimatedTrafficDelay > 0 && (
                            <div className="flex justify-between items-center pt-2 border-t border-gray-200 dark:border-gray-700">
                                <span className="text-gray-900 dark:text-white font-semibold">Total with traffic:</span>
                                <span className="font-bold text-blue-600 dark:text-blue-400 text-lg">
                                    {(() => {
                                        const totalWithTraffic = Math.round((totalDuration + estimatedTrafficDelay) / 60);
                                        const days = Math.floor(totalWithTraffic / (24 * 60));
                                        const hours = Math.floor((totalWithTraffic % (24 * 60)) / 60);
                                        const minutes = totalWithTraffic % 60;
                                        
                                        if (days > 0) {
                                            return `~${days}d ${hours}h`;
                                        } else if (hours > 0) {
                                            return `~${hours}h ${minutes}min`;
                                        } else {
                                            return `~${minutes}min`;
                                        }
                                    })()}
                                </span>
                            </div>
                        )}
                        
                        <div className="flex justify-between items-center pt-3 border-t border-gray-200 dark:border-gray-700">
                            <span className="text-gray-600 dark:text-gray-400">Distance:</span>
                            <span className="font-medium text-gray-900 dark:text-white">
                                {(totalDistance / 1000).toFixed(1)} km
                            </span>
                        </div>
                        
                        <p className="text-xs text-gray-500 dark:text-gray-400 pt-2 border-t border-gray-200 dark:border-gray-700">
                            {estimatedTrafficDelay > 0 
                                ? 'Includes real-time traffic delays' 
                                : 'Includes real-time traffic (no delays detected)'}
                        </p>
                    </div>
                </div>
            )}
        </div>
    );
}