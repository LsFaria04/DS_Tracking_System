export interface TrafficRouteData {
    duration: number;
    durationInTraffic?: number;
    distance: number; 
    trafficDelay: number;
}

export async function getTomTomTrafficData(
    origin: [number, number],
    destination: [number, number],
    apiKey: string
): Promise<TrafficRouteData | null> {
    try {
        const originStr = `${origin[0]},${origin[1]}`;
        const destStr = `${destination[0]},${destination[1]}`;
        
        // TomTom expects lat,lon format
        const url = `https://api.tomtom.com/routing/1/calculateRoute/${originStr}:${destStr}/json?key=${apiKey}&traffic=true&travelMode=car&departAt=now`;

        const response = await fetch(url);
        
        if (!response.ok) {
            console.error('TomTom API error:', response.status);
            return null;
        }

        const data = await response.json();
        
        if (!data.routes?.[0]) {
            return null;
        }

        const route = data.routes[0];
        const summary = route.summary;
        
        return {
            duration: summary.travelTimeInSeconds, // without traffic
            durationInTraffic: summary.trafficDelayInSeconds 
                ? summary.travelTimeInSeconds + summary.trafficDelayInSeconds
                : summary.travelTimeInSeconds,
            distance: summary.lengthInMeters,
            trafficDelay: summary.trafficDelayInSeconds || 0
        };
    } catch (error) {
        console.error('Failed to fetch TomTom traffic data:', error);
        return null;
    }
}
