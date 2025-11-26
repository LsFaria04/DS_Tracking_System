/**
 * Carbon Footprint Calculator for Order Deliveries
 * 
 * Emission factors (kg CO₂e per km):
 * - Road transport (truck): 0.12 kg CO₂e/km
 * - Air freight: 0.50 kg CO₂e/km
 * 
 */

export interface CarbonFootprintData {
    totalEmissions: number;
    roadEmissions: number;
    airEmissions: number;
    distance: number;
    roadDistance: number;
    airDistance: number;
    status: 'CALCULATED' | 'ESTIMATED' | 'ERROR';
    message?: string;
}

// Emission factors (kg CO₂e per km)
const ROAD_EMISSION_FACTOR = 0.12; // Truck/van delivery
const AIR_EMISSION_FACTOR = 0.50; // Air freight/plane

export function calculateCarbonFootprint(
    totalDistanceMeters: number,
    routeSegments: { distance: number; isAir: boolean }[]
): CarbonFootprintData {
    try {
        const totalDistanceKm = totalDistanceMeters / 1000;
        
        let roadDistanceKm = 0;
        let airDistanceKm = 0;
        
        routeSegments.forEach(segment => {
            const segmentKm = segment.distance / 1000;
            if (segment.isAir) {
                airDistanceKm += segmentKm;
            } else {
                roadDistanceKm += segmentKm;
            }
        });
        
        const roadEmissions = roadDistanceKm * ROAD_EMISSION_FACTOR;
        const airEmissions = airDistanceKm * AIR_EMISSION_FACTOR;
        const totalEmissions = roadEmissions + airEmissions;
        
        return {
            totalEmissions,
            roadEmissions,
            airEmissions,
            distance: totalDistanceKm,
            roadDistance: roadDistanceKm,
            airDistance: airDistanceKm,
            status: 'CALCULATED',
        };
    } catch (error) {
        console.error('Failed to calculate carbon footprint:', error);
        return {
            totalEmissions: 0,
            roadEmissions: 0,
            airEmissions: 0,
            distance: 0,
            roadDistance: 0,
            airDistance: 0,
            status: 'ERROR',
            message: 'Unable to calculate carbon footprint',
        };
    }
}

export function getCarbonFootprintExplanation(data: CarbonFootprintData): string {
    if (data.status === 'CALCULATED') {
        return `Based on actual route: ${data.roadDistance.toFixed(1)}km by road (${ROAD_EMISSION_FACTOR} kg CO₂e/km) + ${data.airDistance.toFixed(1)}km by air (${AIR_EMISSION_FACTOR} kg CO₂e/km).`;
    } else if (data.status === 'ESTIMATED') {
        return `Estimated route from seller to delivery location. Road transport emits approximately ${ROAD_EMISSION_FACTOR} kg CO₂e per km, while air freight emits ${AIR_EMISSION_FACTOR} kg CO₂e per km.`;
    } else {
        return 'Unable to calculate carbon footprint at this time.';
    }
}
