"use client";

import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { ReactQueryDevtools } from "@tanstack/react-query-devtools";
import { ReactNode, useState } from "react";

interface ReactQueryProps {
    children: ReactNode;
}

/**
 * ReactQueryProvider component wraps the application with React Query's context.
 * It initializes a QueryClient instance to manage server state and includes the ReactQueryDevtools for debugging.
 *
 * @param {ReactQueryProps} props - Props containing children components to be wrapped.
 * @returns JSX.Element
 */
export default function ReactQueryProvider({ children }: ReactQueryProps) {
    const [queryClient] = useState(() => new QueryClient());

    return (
        <QueryClientProvider client={queryClient}>
            {children}
            <ReactQueryDevtools initialIsOpen={false} />
        </QueryClientProvider>
    );
}
