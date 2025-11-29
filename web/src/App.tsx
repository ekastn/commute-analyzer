import { useState } from "react";
import Sidebar from "./components/Sidebar";
import Map from "./components/Map";
import CommuteForm from "./components/CommuteForm";
import { Plus } from "lucide-react";

function App() {
    const [selectedCommuteId, setSelectedCommuteId] = useState<string | null>(null);
    const [isCreating, setIsCreating] = useState(false);

    return (
        <div className="h-screen flex flex-col md:flex-row">
            <div className="w-full md:w-96 bg-white shadow-xl flex flex-col">
                <div className="p-6 border-b flex items-center justify-between">
                    <h1 className="text-2xl font-bold text-gray-800">CommuteOS</h1>
                    {!isCreating && (
                        <button
                            onClick={() => setIsCreating(true)}
                            className="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-lg flex items-center gap-2 transition"
                        >
                            <Plus size={20} />
                            <span className="hidden sm:inline">Buat Rute Baru</span>
                        </button>
                    )}
                </div>

                <div className="flex-1 overflow-y-auto">
                    {isCreating ? (
                        <CommuteForm onBack={() => setIsCreating(false)} />
                    ) : (
                        <Sidebar onSelect={setSelectedCommuteId} selectedId={selectedCommuteId} />
                    )}
                </div>
            </div>

            {/* Map */}
            <div className="flex-1 relative">
                <Map selectedCommuteId={selectedCommuteId} />
            </div>
        </div>
    );
}

export default App;
