# system-services-ui

A lightweight, secure web-based UI for managing specific **systemd** services on Fedora and other Linux distributions.

[![License: GPL v3](https://img.shields.io/badge/License-GPLv3-blue.svg)](https://www.gnu.org/licenses/gpl-3.0)
[![Fedora Build](https://img.shields.io/badge/Fedora-Ready-blue?logo=fedora)](https://fedoraproject.org/)

---

## 🚀 Features

- **Service Monitoring:** View real-time `active` and `enabled` status of services.
- **Controlled Management:** Restart allowed services directly from a clean web interface.
- **Native Integration:** Built specifically for `systemd` and `polkit`.
- **Zero Runtime Dependencies:** Single Go binary with all static assets (Alpine.js/CSS) embedded via `go:embed`.
- **Enterprise-Ready Packaging:** Includes a production-ready Fedora `.spec` file.

---

## 🛡 Security Model

This tool is built with a **security-first** mindset:

- **Dedicated User:** Runs under its own unprivileged system user (`system-services-ui`).
- **Polkit Authorization:** Does not run as root. Privileged actions are authorized via Polkit rules.
- **Hardened for Fedora:** Ships with Polkit `.rules` and systemd `.service` files out of the box.

---

## 📦 Installation & Deployment

### Fedora RPM (Recommended)

Building and installing the native package ensures all users and permissions are set correctly:

```bash
# Build the archive and SRPM
make archive
rpmbuild -bs packaging/fedora/system-services-ui.spec --define "_sourcedir $(pwd)"

# Install the resulting RPM
sudo dnf install ./system-services-ui-0.1.0-1.fc*.x86_64.rpm
sudo systemctl enable --now system-services-ui
```

---

### Manual Build

```bash
make build
./system-services-ui
```

---

## 🛠 Project Structure

```text
cmd/            Application entry point
internal/       Core logic (API handlers, systemd manager, configuration)
packaging/      Fedora-specific configuration (spec, service, Polkit rules)
web/static/     Frontend assets (HTML, CSS, Alpine.js)
```

---

## ⚙️ Development

The project includes a `Makefile` to automate common tasks:

- `make build` – Compile the Go binary
- `make archive` – Create a vendored source tarball for packaging
- `make clean` – Remove build artifacts

---

## 📄 License

This project is licensed under the **GPL-3.0 License**.  
See the `LICENSE` file for the full text.

---

Created by **Andrei Ivan**
