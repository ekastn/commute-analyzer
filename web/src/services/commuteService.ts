import { api, catchError } from "../lib/api";
import type { Commute } from "../lib/types";

export const listCommutes = (deviceId: string) => {
    return catchError<{ commutes: Commute[]; total: number }>(
        api.get("/commutes", { params: { device_id: deviceId } })
    );
};

export type CreateCommuteInput = Omit<
    Commute,
    | "id"
    | "created_at"
    | "distance_km"
    | "duration_min"
    | "annual_cost_rp"
    | "annual_minutes"
    | "annual_hours"
    | "annual_workdays"
> & { device_id: string };

export const createCommute = (data: CreateCommuteInput) => {
    return catchError<Commute>(api.post("/commutes", data));
};

export const updateCommute = (id: string, data: Partial<Commute>) => {
    return catchError<Commute>(api.patch(`/commutes/${id}`, data));
};

export const deleteCommute = (id: string) => {
    return catchError<void>(api.delete(`/commutes/${id}`));
};
