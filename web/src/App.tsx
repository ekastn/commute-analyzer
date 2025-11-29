import { useState } from "react";
import Sidebar from "./components/Sidebar";
import Map from "./components/Map";
import CommuteForm from "./components/CommuteForm";
import { Plus } from "lucide-react";

function App() {
    const [selectedCommuteId, setSelectedCommuteId] = useState<string | null>(null);
    const [isCreating, setIsCreating] = useState(false);

    const [draftPoints, setDraftPoints] = useState<{
        home: { lat: number; lng: number } | null;
        office: { lat: number; lng: number } | null;
    }>({ home: null, office: null });

    const [pickingMode, setPickingMode] = useState<"home" | "office" | null>(null);

    const handleMapClick = (lat: number, lng: number) => {
        if (!pickingMode) return;
        setDraftPoints((prev) => ({
            ...prev,
            [pickingMode]: { lat, lng },
        }));
        setPickingMode(null);
    };

    return (
        <div className="h-screen flex flex-col md:flex-row">
            <div className="w-full md:w-96 bg-white shadow-xl flex flex-col">
                <div className="p-6 border-b flex items-center justify-between">
                    <h1 className="text-2xl font-bold text-gray-800">CommuteOS</h1>
                    {!isCreating && (
                        <button
                            onClick={() => {
                                setIsCreating(true);
                                setSelectedCommuteId(null);
                                setDraftPoints({ home: null, office: null });
                            }}
                            className="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-lg flex items-center gap-2 transition"
                        >
                            <Plus size={20} />
                            <span className="hidden sm:inline">Buat Rute Baru</span>
                        </button>
                    )}
                </div>

                <div className="flex-1 overflow-y-auto">
                    {isCreating ? (
                        <CommuteForm
                            onBack={() => {
                                setIsCreating(false);
                                setPickingMode(null);
                            }}
                            draftPoints={draftPoints}
                            pickingMode={pickingMode}
                            setPickingMode={setPickingMode}
                        />
                    ) : (
                        <Sidebar onSelect={setSelectedCommuteId} selectedId={selectedCommuteId} />
                    )}
                </div>
            </div>

            <div className="flex-1 relative">
                <Map
                    selectedCommuteId={selectedCommuteId}
                    draftPoints={isCreating ? draftPoints : null}
                    pickingMode={pickingMode}
                    onMapClick={handleMapClick}
                />
            </div>
        </div>
    );
}

export default App;
