import React, { useEffect } from "react";
import { MapContainer, TileLayer, Marker, Popup, useMapEvents, Polyline, useMap } from "react-leaflet";
import { Icon, LatLngBounds } from "leaflet";
import "leaflet/dist/leaflet.css";
import type { Commute } from "../lib/types";

const icon = new Icon({
    iconUrl: "https://unpkg.com/leaflet@1.9.4/dist/images/marker-icon.png",
    shadowUrl: "https://unpkg.com/leaflet@1.9.4/dist/images/marker-shadow.png",
    iconSize: [25, 41],
    iconAnchor: [12, 41],
});

interface MapProps {
    commutes: Commute[];
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

function MapUpdater({ selectedCommute }: { selectedCommute?: Commute }) {
    const map = useMap();

    useEffect(() => {
        if (selectedCommute) {
            const bounds = new LatLngBounds(
                [selectedCommute.home_lat, selectedCommute.home_lng],
                [selectedCommute.office_lat, selectedCommute.office_lng]
            );
            
            // If route geometry exists, include it in bounds for perfect fit
            if (selectedCommute.route_geometry) {
                selectedCommute.route_geometry.forEach(p => {
                    bounds.extend([p[1], p[0]]); // GeoJSON [lng, lat] -> Leaflet [lat, lng]
                });
            }

            map.fitBounds(bounds, { padding: [50, 50] });
        }
    }, [selectedCommute, map]);

    return null;
}

export default function Map({
    commutes,
    selectedCommuteId,
    draftPoints,
    pickingMode,
    onMapClick,
}: MapProps) {
    const center: [number, number] = [-6.2088, 106.8456];
    const selectedCommute = commutes.find((c) => c.id === selectedCommuteId);

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

            <MapUpdater selectedCommute={selectedCommute} />

            {/* Selected Commute */}
            {selectedCommute && (
                <React.Fragment key={selectedCommute.id}>
                    <Marker position={[selectedCommute.home_lat, selectedCommute.home_lng]} icon={icon}>
                        <Popup>Rumah: {selectedCommute.name}</Popup>
                    </Marker>
                    <Marker position={[selectedCommute.office_lat, selectedCommute.office_lng]} icon={icon}>
                        <Popup>Kantor: {selectedCommute.name}</Popup>
                    </Marker>
                    {selectedCommute.route_geometry && (
                        <Polyline
                            positions={selectedCommute.route_geometry.map((p) => [p[1], p[0]] as [number, number])}
                            color="blue"
                            weight={5}
                            opacity={0.8}
                        />
                    )}
                </React.Fragment>
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
