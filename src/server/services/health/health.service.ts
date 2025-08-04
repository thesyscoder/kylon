export class HealthService {
    static ping = () => ({
        status: "ok",
        timestamp: new Date().toISOString(),
    });
}
