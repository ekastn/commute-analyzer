import CommuteCard from "./CommuteCard";
import { Loader2 } from "lucide-react";
import type { Commute } from "../lib/types";

interface SidebarProps {
    commutes: Commute[];
    isLoading: boolean;
    onSelect: (id: string) => void;
    selectedId: string | null;
}

export default function Sidebar({
    commutes,
    isLoading,
    onSelect,
    selectedId,
}: SidebarProps) {
    if (isLoading) {
        return (
            <div className="flex-1 flex items-center justify-center">
                <Loader2 className="animate-spin" size={32} />
            </div>
        );
    }

    if (commutes.length === 0) {
        return (
            <div className="flex-1 flex items-center justify-center p-8 text-center">
                <div>
                    <div className="text-6xl mb-4">Map</div>
                    <p className="text-gray-500">Belum ada rute tersimpan</p>
                    <p className="text-sm text-gray-400 mt-2">Klik "Buat Rute Baru" untuk mulai</p>
                </div>
            </div>
        );
    }

    return (
        <div className="flex-1 overflow-y-auto">
            {commutes.map((commute) => (
                <CommuteCard
                    key={commute.id}
                    commute={commute}
                    isSelected={selectedId === commute.id}
                    onClick={() => onSelect(commute.id)}
                />
            ))}
        </div>
    );
}
