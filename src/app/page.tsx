"use client";

import { trpc } from "@/client";

const Home = () => {
    const heatlth = trpc.health.ping.useQuery();
    if (heatlth.isLoading) return <p>loading...</p>;
    if (heatlth.error) return <p>Error loading health status</p>;
    return (
        <div>
            <h1>Kylon</h1>
            <p>The ultimite K8s distaer recovery</p>
            <h2>Health Status</h2>
            <p>Status:{heatlth.data?.status}</p>
            <p>Timestamp:{heatlth.data?.timestamp}</p>
        </div>
    );
};

export default Home;
