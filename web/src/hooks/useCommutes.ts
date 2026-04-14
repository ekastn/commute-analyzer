import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import {
    listCommutes,
    createCommute,
    updateCommute as updateCommuteService,
    deleteCommute as deleteCommuteService,
} from "../services/commuteService";
import type { Commute } from "../lib/types";
import { useAuth } from "../contexts/AuthContext";

export const useCommutes = () => {
    const queryClient = useQueryClient();
    const { deviceId } = useAuth();

    const { data: commutes = [], isLoading } = useQuery<Commute[]>({
        queryKey: ["commutes"],
        queryFn: async () => {
            if (!deviceId) return [];
            const res = await listCommutes(deviceId);
            return res.commutes;
        },
        enabled: !!deviceId,
    });

    const createMutation = useMutation({
        mutationFn: createCommute,
        onSettled: () => queryClient.invalidateQueries({ queryKey: ["commutes"] }),
    });

    const updateMutation = useMutation({
        mutationFn: ({ id, data }: { id: string; data: Partial<Commute> }) => updateCommuteService(id, data),
        onSettled: () => queryClient.invalidateQueries({ queryKey: ["commutes"] }),
    });

    const deleteMutation = useMutation({
        mutationFn: ({ id }: { id: string }) => deleteCommuteService(id),
        onSettled: () => queryClient.invalidateQueries({ queryKey: ["commutes"] }),
    });

    return {
        commutes,
        isLoading,
        createCommute: createMutation.mutateAsync,
        updateCommute: updateMutation.mutateAsync,
        deleteCommute: deleteMutation.mutateAsync,
    };
};
