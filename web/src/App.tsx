import { useState } from "react";
import Sidebar from "./components/Sidebar";
import Map from "./components/Map";
import CommuteForm from "./components/CommuteForm";
import { Plus, Info } from "lucide-react";
import { useCommutes } from "./hooks/useCommutes";
import AboutModal from "./components/AboutModal";
import type { Commute } from "./lib/types";

function App() {
    const [selectedCommuteId, setSelectedCommuteId] = useState<string | null>(null);
    const [isCreating, setIsCreating] = useState(false);
    const [editingCommute, setEditingCommute] = useState<Commute | null>(null);
    const { commutes, isLoading, updateCommute } = useCommutes();
    const [showAbout, setShowAbout] = useState(false);

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

    const handleEdit = (commute: Commute) => {
        setEditingCommute(commute);
        setSelectedCommuteId(null);
    };

    const handleBack = () => {
        setIsCreating(false);
        setEditingCommute(null);
        setPickingMode(null);
    };

    const isFormOpen = isCreating || !!editingCommute;

    return (
        <div className="h-screen flex flex-col md:flex-row">
            <div className="w-full md:w-96 bg-white shadow-xl flex flex-col">
                <div className="p-6 border-b flex items-center justify-between">
                    <div className="flex items-center gap-2">
                        <h1 className="text-2xl font-bold text-gray-800">Commutes</h1>
                        <button onClick={() => setShowAbout(true)} className="p-1 hover:bg-gray-100 rounded-full transition">
                            <Info size={20} className="text-gray-500" />
                        </button>
                    </div>
                    {!isFormOpen && (
                        <button
                            onClick={() => {
                                setIsCreating(true);
                                setSelectedCommuteId(null);
                                setDraftPoints({ home: null, office: null });
                            }}
                            className="bg-blue-600 hover:bg-blue-700 text-white p-2 rounded-lg transition"
                        >
                            <Plus size={24} />
                        </button>
                    )}
                </div>

                <div className="flex-1 overflow-y-auto">
                    {isFormOpen ? (
                        <CommuteForm
                            onBack={handleBack}
                            draftPoints={draftPoints}
                            pickingMode={pickingMode}
                            setPickingMode={setPickingMode}
                            commute={editingCommute}
                            onEdit={async (data) => {
                                if (editingCommute) {
                                    await updateCommute({ id: editingCommute.id, data });
                                }
                            }}
                        />
                    ) : (
                        <Sidebar
                            commutes={commutes}
                            isLoading={isLoading}
                            onSelect={setSelectedCommuteId}
                            selectedId={selectedCommuteId}
                            onEdit={handleEdit}
                        />
                    )}
                </div>
            </div>

            <div className="flex-1 relative">
                <Map
                    commutes={commutes}
                    selectedCommuteId={selectedCommuteId}
                    draftPoints={isFormOpen ? draftPoints : null}
                    pickingMode={pickingMode}
                    onMapClick={handleMapClick}
                />
            </div>
            {showAbout && <AboutModal onClose={() => setShowAbout(false)} />}
        </div>
    );
}

export default App;
