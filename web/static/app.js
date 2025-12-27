function app() {
    return {
        serverUrl: null,
        connected: false,
        system: {},
        services: [],

        async fetchSystem() {
            try {
                const res = await fetch('/api/system', { cache: 'no-store' });
                this.system = await res.json();
                this.serverUrl = this.system.hostname;
            } catch (e) {
                console.error("System fetch failed", e);
            }
        },

        async fetchServices() {
            try {
                const res = await fetch('/api/services/status', { cache: 'no-store' });
                const data = await res.json();
                this.services = data.services;
            } catch (e) {
                console.error("Services fetch failed", e);
            }
        },

        async fetchAll() {
            try {
                await Promise.all([
                    this.fetchSystem(),
                    this.fetchServices()
                ]);
                this.connected = true;
            } catch {
                this.connected = false;
            }
        },

        async restartService(name) {
            if (!confirm(`Restart ${name}?`)) return;
            try {
                await fetch(`/api/services/${name}/restart`, { method: 'POST' });
                await this.fetchServices();
            } catch (e) {
                alert("Failed to restart service");
            }
        },

        init() {
            this.fetchAll();
            setInterval(() => this.fetchAll(), 5000);
        }
    };
}