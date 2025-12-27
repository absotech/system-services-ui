# system-services-ui

A lightweight web-based UI for managing systemd services on Fedora and similar systems.

## Features
- Service status (enabled/running)
- Restart allowed services
- systemd + polkit integration
- Minimal frontend (Alpine.js)

## Security model
Designed to run on trusted networks (e.g. behind VPN).
Authorization for systemd actions is enforced via polkit.

## License
GPL-3.0
