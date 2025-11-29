import { MapContainer, TileLayer, Marker, Popup } from "react-leaflet";
import { Icon } from "leaflet";
import "leaflet/dist/leaflet.css";

const icon = new Icon({
    iconUrl: "https://unpkg.com/leaflet@1.9.4/dist/images/marker-icon.png",
    shadowUrl: "https://unpkg.com/leaflet@1.9.4/dist/images/marker-shadow.png",
    iconSize: [25, 41],
    iconAnchor: [12, 41],
});

export default function Map({ selectedCommuteId }: { selectedCommuteId: number | null }) {
    const center: [number, number] = [-6.2088, 106.8456];

    return (
        <MapContainer center={center} zoom={11} className="h-full w-full">
            <TileLayer
                url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
                attribution="&copy; OpenStreetMap contributors"
            />
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
        </MapContainer>
    );
}
