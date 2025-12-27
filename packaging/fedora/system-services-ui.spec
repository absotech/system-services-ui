Name:           system-services-ui
Version:        0.1.0
Release:        1%{?dist}
Summary:        Web UI for managing systemd services

License:        GPL-3.0-or-later
URL:            https://github.com/absotech/system-services-ui
Source0:        %{name}-%{version}.tar.gz

BuildRequires:  golang
Requires:       systemd
Requires:       polkit

%description
A lightweight web-based UI for managing selected systemd services
on Fedora systems. Designed for trusted networks.

%pre
getent group system-services-ui >/dev/null || groupadd -r system-services-ui
getent passwd system-services-ui >/dev/null || \
    useradd -r -g system-services-ui -s /sbin/nologin \
    -d / system-services-ui

%prep
%autosetup

%build
export LDFLAGS="-X main.version=%{version} -B 0x$(head -c20 /dev/urandom | od -An -tx1 | tr -d ' \n')"
go build -mod=vendor -v -ldflags "$LDFLAGS" -o %{name} ./cmd/%{name}

%install
install -Dm0755 system-services-ui \
    %{buildroot}%{_bindir}/system-services-ui

install -Dm0644 packaging/fedora/system-services-ui.service \
    %{buildroot}%{_unitdir}/system-services-ui.service

install -Dm0644 packaging/fedora/50-system-services-ui.rules \
    %{buildroot}%{_datadir}/polkit-1/rules.d/50-system-services-ui.rules

%post
%systemd_post system-services-ui.service

%preun
%systemd_preun system-services-ui.service

%postun
%systemd_postun_with_restart system-services-ui.service

%files
%license LICENSE
%doc README.md
%{_bindir}/system-services-ui
%{_unitdir}/system-services-ui.service
%{_datadir}/polkit-1/rules.d/50-system-services-ui.rules

%changelog
* Sat Dec 27 2025 Andrei Ivan <contact@andreiivan.com> - 0.1.0-1
- Initial Fedora package
