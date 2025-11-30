import { X } from "lucide-react";

interface AboutModalProps {
    onClose: () => void;
}

export default function AboutModal({ onClose }: AboutModalProps) {
    return (
        <div className="fixed inset-0 bg-black/50 z-[9999] flex items-center justify-center p-4">
            <div className="bg-white rounded-xl shadow-2xl max-w-md w-full relative overflow-hidden">
                {/* Header */}
                <div className="bg-blue-600 px-6 py-4 flex items-center justify-between">
                    <h2 className="text-xl font-bold text-white">Tentang Aplikasi</h2>
                    <button onClick={onClose} className="text-white hover:bg-blue-700 p-1 rounded-full transition">
                        <X size={24} />
                    </button>
                </div>

                {/* Content */}
                <div className="p-6 space-y-4 text-gray-700">
                    <div>
                        <h3 className="font-bold text-gray-900 mb-2">Apa itu Commute Analyzer?</h3>
                        <p className="text-sm leading-relaxed">
                            Aplikasi ini membantu Anda menghitung estimasi <strong>biaya tahunan</strong> dan <strong>waktu yang terbuang</strong> untuk perjalanan pulang-pergi (commute) ke kantor.
                        </p>
                    </div>

                    <hr className="border-gray-200" />

                    <div>
                        <h3 className="font-bold text-gray-900 mb-2">Cara Penggunaan:</h3>
                        <ul className="text-sm space-y-2 list-decimal list-inside">
                            <li>Klik tombol <span className="inline-block bg-blue-600 text-white rounded px-1 py-0.5 text-xs font-bold">+</span> di pojok kanan atas.</li>
                            <li>Masukkan <strong>Nama Rute</strong> (misal: "Rute Kantor").</li>
                            <li>Klik tombol pilih lokasi <strong>Rumah</strong> & <strong>Kantor</strong>, lalu klik titik di peta.</li>
                            <li>Isi detail kendaraan, harga bensin, dan hari kerja.</li>
                            <li>Klik <strong>Hitung & Simpan</strong> untuk melihat hasilnya.</li>
                        </ul>
                    </div>

                    <div className="bg-blue-50 border border-blue-100 rounded-lg p-3 text-xs text-blue-800">
                        <strong>Tips:</strong> Klik pada kartu rute di sidebar untuk melihat detail jalur perjalanan di peta.
                    </div>
                </div>

                {/* Footer */}
                <div className="bg-gray-50 px-6 py-4 flex justify-end">
                    <button
                        onClick={onClose}
                        className="bg-gray-200 hover:bg-gray-300 text-gray-800 px-4 py-2 rounded-lg text-sm font-medium transition"
                    >
                        Tutup
                    </button>
                </div>
            </div>
        </div>
    );
}
