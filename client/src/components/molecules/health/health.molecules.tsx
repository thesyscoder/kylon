"use client";

import { useHealthCheck } from "@/hooks";
import { useHealthStore } from "@/stores/health/health.store";
import { useEffect } from "react";

const Health = () => {
    const { data, isLoading, isError } = useHealthCheck();
    const healthStatus = useHealthStore((s) => s.healthStatus);
    const setHealthStatus = useHealthStore((s) => s.setHealthStatus);

    useEffect(() => {
        if (isLoading) {
            setHealthStatus("loading...");
        } else if (isError) {
            setHealthStatus("unreachable");
        } else if (data && data.data && data.data.status) {
            // Corrected: Access data.data.status
            setHealthStatus(data.data.status); // Set the status from the nested data object
        }
    }, [data, isError, isLoading, setHealthStatus]);

    return (
        <div>
            <h1>Backend Health</h1>
            <p>Status: {healthStatus}</p>
            {/* Optional: Display more details if needed */}
            {data && data.data && (
                <div>
                    <p>Message: {data.data.message}</p>
                    <p>Application: {data.data.application}</p>
                    <p>Version: {data.data.version}</p>
                    <p>Environment: {data.data.environment}</p>
                    <p>Timestamp: {data.data.timestamp}</p>
                    <h3>Dependencies:</h3>
                    <p>Database: {data.data.dependencies.database}</p>
                    <p>Kubernetes: {data.data.dependencies.kubernetes}</p>
                </div>
            )}
            {isError && (
                <p style={{ color: "red" }}>Error fetching health data.</p>
            )}
        </div>
    );
};

export default Health;
