import axios, { AxiosError, type AxiosResponse } from "axios";

const API_URL = import.meta.env.VITE_API_URL || "http://localhost:8080/api/v1";

export const api = axios.create({
    baseURL: API_URL,
});

export interface APIResponse<T> {
    success: boolean;
    data: T;
    error?: string;
}

export async function catchError<T>(promise: Promise<AxiosResponse<APIResponse<T>>>): Promise<T> {
    try {
        const response = await promise;
        if (!response.data.success) {
            throw new Error(response.data.error || "Unknown API error");
        }
        return response.data.data;
    } catch (error) {
        if (error instanceof AxiosError && error.response?.data) {
            const backendError = error.response.data as APIResponse<null>;
            throw new Error(backendError.error || error.message);
        }
        throw error;
    }
}
