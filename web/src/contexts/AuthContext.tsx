import { createContext, useContext, useState, useEffect, useCallback, type ReactNode } from "react";
import { v4 as uuidv4 } from "uuid";

interface AuthContextType {
    deviceId: string | null;
    clearAuth: () => void;
    isAuthenticated: () => boolean;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function AuthProvider({ children }: { children: ReactNode }) {
    const [deviceId, setDeviceId] = useState<string | null>(null);

    const isAuthenticated = () => !!deviceId;

    useEffect(() => {
        let storedDeviceId = localStorage.getItem("device_id");
        if (!storedDeviceId) {
            storedDeviceId = uuidv4();
            localStorage.setItem("device_id", storedDeviceId!);
        }
        setDeviceId(storedDeviceId);
    }, []);

    const clearAuth = useCallback(() => {
        localStorage.removeItem("jwt_token");
    }, []);

    return (
        <AuthContext.Provider value={{ deviceId, clearAuth, isAuthenticated }}>
            {children}
        </AuthContext.Provider>
    );
}

export function useAuth() {
    const context = useContext(AuthContext);
    if (context === undefined) {
        throw new Error("useAuth must be used within an AuthProvider");
    }
    return context;
}
