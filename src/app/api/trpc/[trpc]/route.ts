import { fetchRequestHandler } from "@trpc/server/adapters/fetch";
import { appRouter } from "@/server/api";
import { createContext } from "@/server/trpc";

const handler = (req: Request) => {
    console.log("[tRPC] Incoming request", req.url);
    return fetchRequestHandler({
        endpoint: "/api/trpc",
        req,
        router: appRouter,
        createContext,
        onError({ error }) {
            console.error("[tRPC Error]", error);
        },
    });
};

export { handler as GET, handler as POST };
