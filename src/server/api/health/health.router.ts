import { HealthService } from "@/server/services";
import { publicProcedure, router } from "@/server/trpc";

export class HealthRouter {
    static create = () =>
        router({
            ping: publicProcedure.query(() => HealthService.ping()),
        });
}
