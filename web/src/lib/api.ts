import axios from "axios";

const API_URL = import.meta.env.VITE_API_URL || "http://localhost:8080/api/v1";

const api = axios.create({
    baseURL: API_URL,
});

export const commutesApi = {
    create: (data: any) => api.post("/commutes", data),
    list: (deviceId: string) => api.get("/commutes", { params: { device_id: deviceId } }),
    update: (id: number, data: any) => api.patch(`/commutes/${id}`, data),
    delete: (id: number) => api.delete(`/commutes/${id}`),
};
