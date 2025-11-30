import { useState } from "react";
import { useCommutes } from "../hooks/useCommutes";
import { ArrowLeft, MapPin, Fuel, Calendar } from "lucide-react";
import { useAuth } from "../contexts/AuthContext";

interface CommuteFormProps {
    onBack: () => void;
    draftPoints: {
        home: { lat: number; lng: number } | null;
        office: { lat: number; lng: number } | null;
    };
    pickingMode: "home" | "office" | null;
    setPickingMode: (mode: "home" | "office" | null) => void;
}

export default function CommuteForm({
    onBack,
    draftPoints,
    pickingMode,
    setPickingMode,
}: CommuteFormProps) {
    const { createCommute } = useCommutes();
    const { deviceId } = useAuth();
    const [loading, setLoading] = useState(false);
    const [name, setName] = useState("");

    const [form, setForm] = useState({
        vehicle: "motorcycle" as "car" | "motorcycle",
        fuel_price: 10000,
        days_per_week: 5,
    });

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        if (!draftPoints.home || !draftPoints.office) {
            alert("Mohon pilih lokasi Rumah dan Kantor di peta.");
            return;
        }

        if (!name.trim()) {
            alert("Mohon masukkan nama rute.");
            return;
        }

        setLoading(true);
        try {
            await createCommute({
                device_id: deviceId!,
                name: name.trim(),
                home_lat: draftPoints.home.lat,
                home_lng: draftPoints.home.lng,
                office_lat: draftPoints.office.lat,
                office_lng: draftPoints.office.lng,
                ...form,
            });
            onBack();
        } catch (err) {
            console.error("Save failed:", err);
            alert("Gagal menyimpan. Pastikan backend aktif.");
        } finally {
            setLoading(false);
        }
    };

    return (
        <div className="h-full flex flex-col">
            {/* Header */}
            <div className="bg-white border-b px-6 py-4 flex items-center gap-3">
                <button onClick={onBack} className="p-2 hover:bg-gray-100 rounded-lg">
                    <ArrowLeft size={24} />
                </button>
                <h2 className="text-xl font-bold">Buat Rute Baru</h2>
            </div>

            <div className="flex-1 overflow-y-auto p-6 space-y-6">
                <div className="bg-blue-50 border border-blue-200 rounded-lg p-4 text-sm">
                    Pilih lokasi Rumah dan Kantor dengan mengklik tombol di bawah ini, lalu klik pada peta.
                </div>

                {/* Commute Name Input */}
                <div>
                    <label htmlFor="commute-name" className="block text-sm font-medium mb-2">
                        Nama Rute
                    </label>
                    <input
                        id="commute-name"
                        type="text"
                        value={name}
                        onChange={(e) => setName(e.target.value)}
                        className="w-full px-4 py-3 border rounded-lg focus:ring-2 focus:ring-blue-500"
                        placeholder="Ex: Rute Harian Kantor"
                        required
                    />
                </div>

                {/* Location Picker Buttons */}
                <div className="space-y-4">
                    <button
                        type="button"
                        onClick={() => setPickingMode("home")}
                        className={`w-full rounded-lg shadow-sm border p-4 text-left transition-all ${
                            pickingMode === "home"
                                ? "ring-2 ring-blue-500 bg-blue-50 text-blue-800 border-transparent"
                                : "bg-white hover:border-blue-300"
                        }`}
                    >
                        <div className="flex items-center gap-3">
                            <div className="bg-green-100 p-2 rounded-full">
                                <MapPin className="text-green-600" size={20} />
                            </div>
                            <div>
                                <p className="font-medium">Lokasi Rumah</p>
                                <p className="text-sm text-gray-600 truncate">
                                    {draftPoints.home
                                        ? `${draftPoints.home.lat.toFixed(5)}, ${draftPoints.home.lng.toFixed(5)}`
                                        : "Klik untuk pilih di peta"}
                                </p>
                            </div>
                        </div>
                    </button>

                    <button
                        type="button"
                        onClick={() => setPickingMode("office")}
                        className={`w-full rounded-lg shadow-sm border p-4 text-left transition-all ${
                            pickingMode === "office"
                                ? "ring-2 ring-blue-500 bg-blue-50 text-blue-800 border-transparent"
                                : "bg-white hover:border-blue-300"
                        }`}
                    >
                        <div className="flex items-center gap-3">
                            <div className="bg-red-100 p-2 rounded-full">
                                <MapPin className="text-red-600" size={20} />
                            </div>
                            <div>
                                <p className="font-medium">Lokasi Kantor</p>
                                <p className="text-sm text-gray-600 truncate">
                                    {draftPoints.office
                                        ? `${draftPoints.office.lat.toFixed(5)}, ${draftPoints.office.lng.toFixed(5)}`
                                        : "Klik untuk pilih di peta"}
                                </p>
                            </div>
                        </div>
                    </button>
                </div>

                {/* Form */}
                <form onSubmit={handleSubmit} className="space-y-5">
                    <div>
                        <label className="block text-sm font-medium mb-2">
                            <Fuel className="inline mr-2" size={18} />
                            Kendaraan
                        </label>
                        <select
                            value={form.vehicle}
                            onChange={(e) => setForm({ ...form, vehicle: e.target.value as any })}
                            className="w-full px-4 py-3 border rounded-lg focus:ring-2 focus:ring-blue-500"
                        >
                            <option value="motorcycle">Motor</option>
                            <option value="car">Mobil</option>
                        </select>
                    </div>

                    <div>
                        <label className="block text-sm font-medium mb-2">
                            Harga Bensin / Liter
                        </label>
                        <input
                            type="number"
                            value={form.fuel_price}
                            onChange={(e) =>
                                setForm({ ...form, fuel_price: Number(e.target.value) })
                            }
                            className="w-full px-4 py-3 border rounded-lg focus:ring-2 focus:ring-blue-500"
                            required
                        />
                    </div>

                    <div>
                        <label className="block text-sm font-medium mb-2">
                            <Calendar className="inline mr-2" size={18} />
                            Hari Masuk Kantor / Minggu
                        </label>
                        <div className="grid grid-cols-7 gap-2">
                            {[1, 2, 3, 4, 5, 6, 7].map((d) => (
                                <button
                                    key={d}
                                    type="button"
                                    onClick={() => setForm({ ...form, days_per_week: d })}
                                    className={`py-3 rounded-lg font-medium transition ${
                                        form.days_per_week === d
                                            ? "bg-blue-600 text-white"
                                            : "bg-gray-100 hover:bg-gray-200"
                                    }`}
                                >
                                    {d}
                                </button>
                            ))}
                        </div>
                    </div>

                    <button
                        type="submit"
                        disabled={loading || !draftPoints.home || !draftPoints.office}
                        className="w-full bg-blue-600 hover:bg-blue-700 disabled:bg-gray-400 text-white font-bold py-4 rounded-lg transition"
                    >
                        {loading ? "Menyimpan..." : "Hitung & Simpan"}
                    </button>
                </form>
            </div>
        </div>
    );
}
