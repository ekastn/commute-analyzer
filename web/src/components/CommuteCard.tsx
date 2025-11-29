import type { Commute } from "../lib/types";
import { Trash2, Car, Bike } from "lucide-react";
import { useCommutes } from "../hooks/useCommutes";

export default function CommuteCard({
    commute,
    isSelected,
    onClick,
}: {
    commute: Commute;
    isSelected: boolean;
    onClick: () => void;
}) {
    const { deleteCommute } = useCommutes();

    return (
        <div
            className={`p-4 border-b hover:bg-gray-50 cursor-pointer transition-all ${isSelected ? "bg-blue-50 border-l-4 border-l-blue-600" : ""}`}
            onClick={onClick}
        >
            <div className="flex items-start justify-between">
                <div className="flex-1">
                    <h3 className="font-semibold text-gray-800">{commute.name || "Rute Baru"}</h3>
                    <div className="flex items-center gap-4 mt-2 text-sm text-gray-600">
                        <span className="flex items-center gap-1">
                            {commute.vehicle === "car" ? <Car size={16} /> : <Bike size={16} />}
                            {commute.distance_km.toFixed(1)} km · {commute.duration_min.toFixed(0)}{" "}
                            menit
                        </span>
                    </div>
                    <div className="mt-3 space-y-1">
                        <p className="text-lg font-bold text-red-600">
                            Rp {(commute.annual_cost_rp / 1_000_000).toFixed(2)} jt/tahun
                        </p>
                        <p className="text-sm text-gray-600">
                            {commute.annual_workdays.toFixed(0)} hari kerja hilang di jalan
                        </p>
                    </div>
                </div>
                <button
                    onClick={(e) => {
                        e.stopPropagation();
                        deleteCommute({ id: commute.id });
                    }}
                    className="text-gray-400 hover:text-red-600 transition"
                >
                    <Trash2 size={18} />
                </button>
            </div>
        </div>
    );
}
