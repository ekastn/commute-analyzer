import { MapContainer, TileLayer, Marker, Popup, useMapEvents } from "react-leaflet";
import { Icon } from "leaflet";
import "leaflet/dist/leaflet.css";

const icon = new Icon({
    iconUrl: "https://unpkg.com/leaflet@1.9.4/dist/images/marker-icon.png",
    shadowUrl: "https://unpkg.com/leaflet@1.9.4/dist/images/marker-shadow.png",
    iconSize: [25, 41],
    iconAnchor: [12, 41],
});

interface MapProps {
    selectedCommuteId: string | null;
    draftPoints?: {
        home: { lat: number; lng: number } | null;
        office: { lat: number; lng: number } | null;
    } | null;
    pickingMode?: "home" | "office" | null;
    onMapClick?: (lat: number, lng: number) => void;
}

function MapEvents({
    isActive,
    onClick,
}: {
    isActive: boolean;
    onClick: (lat: number, lng: number) => void;
}) {
    useMapEvents({
        click(e) {
            if (isActive) {
                onClick(e.latlng.lat, e.latlng.lng);
            }
        },
    });
    return null;
}

export default function Map({
    selectedCommuteId,
    draftPoints,
    pickingMode,
    onMapClick,
}: MapProps) {
    const center: [number, number] = [-6.2088, 106.8456];

    return (
        <MapContainer
            center={center}
            zoom={11}
            className={`h-full w-full ${pickingMode ? "cursor-crosshair" : ""}`}
        >
            <TileLayer
                url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
                attribution="&copy; OpenStreetMap contributors"
            />
            
            {onMapClick && (
                <MapEvents isActive={!!pickingMode} onClick={onMapClick} />
            )}

            {/* Existing Commutes (Placeholder logic for now) */}
            {selectedCommuteId && (
                <>
                    <Marker position={[-6.2088, 106.8456]} icon={icon}>
                        <Popup>Rumah</Popup>
                    </Marker>
                    <Marker position={[-6.1892, 106.7891]} icon={icon}>
                        <Popup>Kantor</Popup>
                    </Marker>
                </>
            )}

            {/* Draft Points */}
            {draftPoints?.home && (
                <Marker position={[draftPoints.home.lat, draftPoints.home.lng]} icon={icon}>
                    <Popup>Rumah (Baru)</Popup>
                </Marker>
            )}
            {draftPoints?.office && (
                <Marker position={[draftPoints.office.lat, draftPoints.office.lng]} icon={icon}>
                    <Popup>Kantor (Baru)</Popup>
                </Marker>
            )}
        </MapContainer>
    );
}
