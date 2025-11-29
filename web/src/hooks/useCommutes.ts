import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import {
    listCommutes,
    createCommute,
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
        onSuccess: () => queryClient.invalidateQueries({ queryKey: ["commutes"] }),
    });

    const deleteMutation = useMutation({
        mutationFn: ({ id }: { id: string }) => deleteCommuteService(id),
        onSuccess: () => queryClient.invalidateQueries({ queryKey: ["commutes"] }),
    });

    return {
        commutes,
        isLoading,
        createCommute: createMutation.mutateAsync,
        deleteCommute: deleteMutation.mutateAsync,
    };
};
