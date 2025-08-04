import { router } from "../trpc";
import { HealthRouter } from "./health/health.router";

export const appRouter = router({
    health: HealthRouter.create(),
});

export type AppRouter = typeof appRouter;
