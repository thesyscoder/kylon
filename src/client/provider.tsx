"use client";

import { httpBatchLink, loggerLink } from "@trpc/client";
import { trpc } from "./index";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { PropsWithChildren, useState } from "react";

export const TrpcProvider = ({ children }: PropsWithChildren) => {
    const [queryClient] = useState(() => new QueryClient());

    const [trpcClient] = useState(() =>
        trpc.createClient({
            links: [
                loggerLink({
                    enabled: () => process.env.NODE_ENV === "development",
                }),
                httpBatchLink({
                    url: "/api/trpc",
                }),
            ],
        }),
    );

    return (
        <trpc.Provider client={trpcClient} queryClient={queryClient}>
            <QueryClientProvider client={queryClient}>
                {children}
            </QueryClientProvider>
        </trpc.Provider>
    );
};
