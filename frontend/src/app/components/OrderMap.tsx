"use client";

import { MapContainer, TileLayer, Marker, Popup, Polyline, useMap, useMapEvents } from 'react-leaflet';
import { OrderStatus } from '@/app/types';
import 'leaflet/dist/leaflet.css';
import L from 'leaflet';
import { useEffect, useState } from 'react';

// Fix for default marker icons in Next.js
delete (L.Icon.Default.prototype as any)._getIconUrl;
L.Icon.Default.mergeOptions({
    iconRetinaUrl: 'https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.7.1/images/marker-icon-2x.png',
    iconUrl: 'https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.7.1/images/marker-icon.png',
    shadowUrl: 'https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.7.1/images/marker-shadow.png',
});

// Create a custom icon for delivery destination
const deliveryIcon = L.icon({
    iconUrl: 'https://raw.githubusercontent.com/pointhi/leaflet-color-markers/master/img/marker-icon-2x-green.png',
    shadowUrl: 'https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.7.1/images/marker-shadow.png',
    iconSize: [25, 41],
    iconAnchor: [12, 41],
    popupAnchor: [1, -34],
    shadowSize: [41, 41]
});

interface OrderMapProps {
    orderHistory: OrderStatus[];
    deliveryAddress?: string;
}

interface DeliveryCoords {
    lat: number;
    lon: number;
}

// Component to adjust line weight based on zoom
function DynamicPolylines({ routeCoordinates, deliveryCoords }: { 
    routeCoordinates: [number, number][]; 
    deliveryCoords: DeliveryCoords | null 
}) {
    const [zoom, setZoom] = useState(7);
    
    useMapEvents({
        zoomend: (e) => {
            setZoom(e.target.getZoom());
        }
    });

    // Adjust weight based on zoom level
    const getLineWeight = () => {
        if (zoom < 8) return 3;
        if (zoom < 10) return 2.5;
        if (zoom < 12) return 2;
        if (zoom < 14) return 1.5;
        return 1;
    };

    const weight = getLineWeight();

    return (
        <>
            {/* Draw completed route (blue solid line) */}
            {routeCoordinates.length > 1 && (
                <Polyline
                    positions={routeCoordinates}
                    color="blue"
                    weight={weight}
                    opacity={0.7}
                />
            )}

            {/* Draw remaining route to delivery (dashed gray line) */}
            {deliveryCoords && routeCoordinates.length > 0 && (
                <Polyline
                    positions={[
                        routeCoordinates[routeCoordinates.length - 1],
                        [deliveryCoords.lat, deliveryCoords.lon]
                    ]}
                    color="gray"
                    weight={weight}
                    opacity={0.5}
                    dashArray="10, 10"
                />
            )}
        </>
    );
}

export default function OrderMap({ orderHistory, deliveryAddress }: OrderMapProps) {
    const [deliveryCoords, setDeliveryCoords] = useState<DeliveryCoords | null>(null);
    const [geocoding, setGeocoding] = useState(false);
    const [mapInstance, setMapInstance] = useState<L.Map | null>(null);

    const locationsWithCoords = orderHistory.filter(h => h.Storage);

    const routeCoordinates: [number, number][] = locationsWithCoords.map(h => [
        h.Storage!.Latitude,
        h.Storage!.Longitude
    ]);

    //////////////////////////////////////////////////////////////////////////////////////////////
    // Better approach we should do is to geocode the address on the backend when order is created
    //////////////////////////////////////////////////////////////////////////////////////////////
    // Geocode delivery address
    useEffect(() => {
        if (deliveryAddress && !deliveryCoords && !geocoding) {
            setGeocoding(true);
            
            // Nominatim API for geocoding
            fetch(`https://nominatim.openstreetmap.org/search?format=json&q=${encodeURIComponent(deliveryAddress)}&limit=1`)
                .then(res => res.json())
                .then(data => {
                    if (data && data.length > 0) {
                        setDeliveryCoords({
                            lat: parseFloat(data[0].lat),
                            lon: parseFloat(data[0].lon)
                        });
                    }
                })
                .catch(err => console.error('Geocoding error:', err))
                .finally(() => setGeocoding(false));
        }
    }, [deliveryAddress, deliveryCoords, geocoding]);

    // Add delivery coordinates to route if available
    const fullRoute = deliveryCoords 
        ? [...routeCoordinates, [deliveryCoords.lat, deliveryCoords.lon] as [number, number]]
        : routeCoordinates;

    // Calculate center including delivery address
    const center: [number, number] = fullRoute.length > 0
        ? [
            fullRoute.reduce((sum, coord) => sum + coord[0], 0) / fullRoute.length,
            fullRoute.reduce((sum, coord) => sum + coord[1], 0) / fullRoute.length
        ]
        : [39.5, -8.0]; // Default center (Portugal)

    if (locationsWithCoords.length === 0) {
        return (
            <div className="w-full h-96 bg-gray-200 rounded-lg flex items-center justify-center">
                <p className="text-gray-500">No location data available</p>
            </div>
        );
    }

    // Get the last location (current/most recent position)
    const lastLocation = locationsWithCoords[locationsWithCoords.length - 1];

    const handleZoomToDelivery = () => {
        if (deliveryCoords && mapInstance) {
            mapInstance.flyTo([deliveryCoords.lat, deliveryCoords.lon], 15, {
                duration: 1.5
            });
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

            {/* Dynamic polylines that adjust weight based on zoom */}
            <DynamicPolylines 
                routeCoordinates={routeCoordinates} 
                deliveryCoords={deliveryCoords} 
            />

            {/* Add markers for each storage location */}
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

            {/* Add delivery destination marker */}
            {deliveryCoords && (
                <Marker
                    position={[deliveryCoords.lat, deliveryCoords.lon]}
                    icon={deliveryIcon}
                >
                    <Popup>
                        <div className="text-sm">
                            <p className="font-bold text-green-600">Delivery Destination</p>
                            <p className="text-xs mt-1">{deliveryAddress}</p>
                        </div>
                    </Popup>
                </Marker>
            )}

            {/* Show info overlay */}
            {deliveryAddress && (
                <div className="leaflet-bottom leaflet-right" style={{ pointerEvents: 'none' }}>
                    <div 
                        className="bg-white p-3 m-4 rounded-lg shadow-lg cursor-pointer hover:shadow-xl transition-shadow" 
                        style={{ pointerEvents: 'auto' }}
                        onClick={handleZoomToDelivery}
                    >
                        <p className="text-xs font-semibold text-gray-700">Delivery To:</p>
                        <p className="text-xs text-gray-600">{deliveryAddress}</p>
                        <p className="text-xs text-blue-600 mt-1">
                            Current: {lastLocation.Storage!.Name}
                        </p>
                        {geocoding && (
                            <p className="text-xs text-gray-400 mt-1">Locating address...</p>
                        )}
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