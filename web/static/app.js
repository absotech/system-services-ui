function app() {
    return {
        serverUrl: null,
        connected: false,
        loading: true,
        refreshing: false,
        system: {},
        services: [],
        lastUpdated: null,
        restartingService: null,
        toast: null,
        toastTimeout: null,
        refreshTimer: null,

        get runningServices() {
            return this.services.filter(service => service.running).length;
        },

        get loadAverage() {
            const match = (this.system.load || '').match(/load averages?:\s*([\d.]+)\s*,\s*([\d.]+)\s*,\s*([\d.]+)/i);
            return match ? match.slice(1).join(' · ') : '—';
        },

        get connectionLabel() {
            if (this.loading) return 'Connecting';
            return this.connected ? 'Connected' : 'Disconnected';
        },

        async fetchSystem() {
            const res = await fetch('/api/v1/system', { cache: 'no-store' });
            if (!res.ok) throw new Error(`System request failed (${res.status})`);
            this.system = await res.json();
            this.serverUrl = this.system.hostname;
        },

        async fetchServices() {
            const res = await fetch('/api/v1/services/status', { cache: 'no-store' });
            if (!res.ok) throw new Error(`Services request failed (${res.status})`);
            const data = await res.json();
            this.services = Object.entries(data.services || {})
                .map(([name, status]) => ({
                    name,
                    enabled: Boolean(status.enabled),
                    running: Boolean(status.running)
                }))
                .sort((a, b) => a.name.localeCompare(b.name));
        },

        async fetchAll() {
            if (this.refreshing) return;
            this.refreshing = true;
            try {
                await Promise.all([
                    this.fetchSystem(),
                    this.fetchServices()
                ]);
                this.connected = true;
                this.lastUpdated = new Date();
            } catch (e) {
                console.error('Dashboard refresh failed', e);
                this.connected = false;
            } finally {
                this.loading = false;
                this.refreshing = false;
            }
        },

        async restartService(name) {
            if (!confirm(`Restart ${name}?`)) return;
            this.restartingService = name;
            try {
                const res = await fetch(`/api/v1/services/${encodeURIComponent(name)}/restart`, { method: 'POST' });
                if (!res.ok) throw new Error(`Restart request failed (${res.status})`);
                await this.fetchServices();
                this.showToast(`${name} restarted successfully`, 'success');
            } catch (e) {
                console.error(`Failed to restart ${name}`, e);
                this.showToast(`Could not restart ${name}`, 'error');
            } finally {
                this.restartingService = null;
            }
        },

        showToast(message, type) {
            clearTimeout(this.toastTimeout);
            this.toast = { message, type };
            this.toastTimeout = setTimeout(() => {
                this.toast = null;
            }, 3500);
        },

        init() {
            this.fetchAll();
            this.refreshTimer = setInterval(() => this.fetchAll(), 5000);
        },

        destroy() {
            clearInterval(this.refreshTimer);
            clearTimeout(this.toastTimeout);
        }
    };
}