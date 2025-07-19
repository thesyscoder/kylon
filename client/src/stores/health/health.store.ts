import { create } from "zustand";
interface HealthStore {
    healthStatus: string;
    setHealthStatus: (status: string) => void;
}

export const useHealthStore = create<HealthStore>((set) => ({
    healthStatus: "loading...",
    setHealthStatus: (status) => set({ healthStatus: status }),
}));
