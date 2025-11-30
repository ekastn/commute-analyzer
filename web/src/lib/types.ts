export interface Commute {
    id: string;
    name: string;
    home_lng: number;
    home_lat: number;
    office_lat: number;
    office_lng: number;
    route_geometry?: number[][];
    distance_km: number;
    duration_min: number;
    vehicle: "car" | "motorcycle";
    fuel_price: number;
    days_per_week: number;
    annual_cost_rp: number;
    annual_minutes: number;
    annual_hours: number;
    annual_workdays: number;
    created_at: string;
}
