import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { commutesApi } from "../lib/api";
import type { Commute } from "../lib/types";
import { useAuth } from "../contexts/AuthContext";

export const useCommutes = () => {
    const queryClient = useQueryClient();
    const { deviceId } = useAuth();

    const { data: commutes = [], isLoading } = useQuery<Commute[]>({
        queryKey: ["commutes"],
        queryFn: async () => {
            const res = await commutesApi.list(deviceId!);
            return res.data.commutes;
        },
    });

    const createMutation = useMutation({
        mutationFn: commutesApi.create,
        onSuccess: () => queryClient.invalidateQueries({ queryKey: ["commutes"] }),
    });

    const deleteMutation = useMutation({
        mutationFn: commutesApi.delete,
        onSuccess: () => queryClient.invalidateQueries({ queryKey: ["commutes"] }),
    });

    return {
        commutes,
        isLoading,
        createCommute: createMutation.mutateAsync,
        deleteCommute: deleteMutation.mutateAsync,
    };
};
