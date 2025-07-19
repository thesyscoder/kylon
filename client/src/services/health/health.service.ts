import { axiosClient } from "@/libs";
import { HealthzResponse } from "@/types/health/health.types";
import { handleError } from "@/utils";
export class HealthCheckService {
    async getHealthCheckStatus(): Promise<HealthzResponse> {
        try {
            const response = await axiosClient
                .getInstance()
                .get<HealthzResponse>("/healthz");
            return response.data;
        } catch (error) {
            handleError("[HealtService Error]:", error);
        }
    }
}
