export interface DependenciesStatus {
    database: string;
    kubernetes: string;
}

export interface HealthzData {
    status: string;
    message: string;
    application: string;
    version: string;
    environment: string;
    timestamp: string;
    dependencies: DependenciesStatus;
}

export interface HealthzResponse {
    success: boolean;
    message: string;
    data: HealthzData;
}
