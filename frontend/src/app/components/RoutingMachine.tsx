import { useEffect, useRef } from 'react';
import { useMap } from 'react-leaflet';
import L from 'leaflet';
import 'leaflet-routing-machine';

interface RoutingProps {
    waypoints: [number, number][];
    routeColor: 'green' | 'red' | 'blue';
}

// Icon fix (unchanged)
delete (L.Icon.Default.prototype as any)._getIconUrl;
L.Icon.Default.mergeOptions({
    iconRetinaUrl: 'https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.7.1/images/marker-icon-2x.png',
    iconUrl: 'https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.7.1/images/marker-icon.png',
    shadowUrl: 'https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.7.1/images/marker-shadow.png',
});


export default function RoutingMachine({ waypoints, routeColor }: RoutingProps) {
    const map = useMap();
    const routingRef = useRef<L.Routing.Control | null>(null);


    useEffect(() => {
        if (!map || waypoints.length < 2) return;

        const leafletWaypoints = waypoints.map(coords => L.latLng(coords[0], coords[1]));
        
        const lineOptions: L.Routing.LineOptions = {
            styles: [
                { 
                    color: routeColor,
                    opacity: 0.7, 
                    weight: 5
                }
            ],
            extendToWaypoints: false,
            missingRouteTolerance: 5
        };

        const routingControl = L.Routing.control({
            waypoints: leafletWaypoints,
            routeWhileDragging: false,
            addWaypoints: false,
            show: false,
            lineOptions: lineOptions,
            fitSelectedRoutes: false,
            router: L.Routing.osrmv1({
                serviceUrl: 'https://routing.openstreetmap.de/routed-car/route/v1'
            }),
            // @ts-expect-error: createMarker is supported by leaflet-routing-machine
            createMarker: () => { 
                return null; 
            }
        }).addTo(map);

        routingRef.current = routingControl;

        return () => {
            try {
                (routingControl as any).setWaypoints([]);
            } catch (e) {}
        
            try {
                map.removeControl(routingControl);
            } catch (e) {}
        };

    }, [map, waypoints, routeColor]);

    return null;
}