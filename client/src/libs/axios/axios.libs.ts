import axios, { AxiosInstance } from "axios";

class AxiosClient {
    private instance: AxiosInstance;
    constructor() {
        this.instance = axios.create({
            baseURL: process.env.NEXT_PUBLIC_API_BASE_URL || "/api", // Use "/api" for Next.js rewrites
            timeout: 10000,
            withCredentials: true,
            headers: {
                "Content-Type": "application/json",
            },
        });
        // TODO: Add request and response interceptors
    }

    public getInstance(): AxiosInstance {
        return this.instance;
    }
}

export const axiosClient = new AxiosClient();
