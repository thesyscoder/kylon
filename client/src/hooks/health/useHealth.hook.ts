import { HealthCheckService } from "@/services";
import { useQuery } from "@tanstack/react-query";

const healthcheckService = new HealthCheckService();

export const useHealthCheck = () => {
    return useQuery({
        queryKey: ["healthcheck"],
        staleTime: 5 * 60 * 1000,
        gcTime: 10 * 60 * 1000,
        retry: false,
        refetchOnWindowFocus: false,
        queryFn: async () => {
            try {
                const response =
                    await healthcheckService.getHealthCheckStatus();
                return response;
            } catch (error) {
                throw error;
            }
        },
    });
};
