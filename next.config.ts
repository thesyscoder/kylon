// next.config.js

import type { NextConfig } from "next";
import path from "path";

const nextConfig: NextConfig = {
    sassOptions: {
        includePaths: [path.join(__dirname, "src", "styles")],
    },
    async rewrites() {
        return [
            {
                source: "/api/:path*",
                destination: "http://localhost:8080/api/:path*",
            },
        ];
    },
};

export default nextConfig;
