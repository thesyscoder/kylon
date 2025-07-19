/**
 * handleError.ts
 * A utility to standardize error handling across services.
 */

export default function handleError(context: string, error: unknown): never {
    const message =
        error instanceof Error
            ? error.message
            : `${context} failed in service.`;

    console.error(`Error in ${context}:`, message);
    throw new Error(message);
}
