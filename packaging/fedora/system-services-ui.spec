%bcond check 1

Name:           system-services-ui
Version:        0.1.0
Release:        1%{?dist}
Summary:        Web UI for managing systemd services

# system-services-ui itself: GPL-3.0-or-later
# Bundled Go dependencies: Apache-2.0, BSD-3-Clause, MIT
# Verify with: go_vendor_license --config go-vendor-tools.toml report expression
License:        GPL-3.0-or-later AND Apache-2.0 AND BSD-3-Clause AND MIT
URL:            https://github.com/absotech/system-services-ui
Source0:        %{url}/archive/v%{version}/%{name}-%{version}.tar.gz
# Generated with: go_vendor_archive create --config go-vendor-tools.toml system-services-ui.spec
Source1:        %{name}-%{version}-vendor.tar.bz2
Source2:        go-vendor-tools.toml

ExclusiveArch:  %{golang_arches_future}

BuildRequires:  go-rpm-macros
BuildRequires:  go-vendor-tools
BuildRequires:  systemd-rpm-macros
Requires:       systemd
Requires:       polkit

%description
A lightweight web-based UI for managing selected systemd services
on Fedora systems. Designed for trusted networks.

%prep
%autosetup -p1
tar -xf %{SOURCE1}

%generate_buildrequires
%go_vendor_license_buildrequires -c %{SOURCE2}

%build
%global gomodulesmode GO111MODULE=on
%gobuild -o %{name} ./cmd/%{name}

%install
%go_vendor_license_install -c %{SOURCE2}

install -Dm0755 system-services-ui \
    %{buildroot}%{_bindir}/system-services-ui

install -Dm0644 packaging/fedora/system-services-ui.service \
    %{buildroot}%{_unitdir}/system-services-ui.service

install -Dm0644 packaging/fedora/50-system-services-ui.rules \
    %{buildroot}%{_datadir}/polkit-1/rules.d/50-system-services-ui.rules

install -Dm0644 internal/config/system-services-ui.yaml \
    %{buildroot}%{_sysconfdir}/system-services-ui/system-services-ui.yaml

install -Dm0644 packaging/fedora/system-services-ui.conf \
    %{buildroot}%{_sysusersdir}/system-services-ui.conf

%check
%go_vendor_license_check -c %{SOURCE2}
%if %{with check}
%gotest ./...
%endif

%post
%systemd_post system-services-ui.service

%preun
%systemd_preun system-services-ui.service

%postun
%systemd_postun_with_restart system-services-ui.service

%files -f %{go_vendor_license_filelist}
%doc README.md
%{_bindir}/system-services-ui
%{_unitdir}/system-services-ui.service
%{_datadir}/polkit-1/rules.d/50-system-services-ui.rules
%{_sysusersdir}/system-services-ui.conf
%dir %{_sysconfdir}/system-services-ui
%config(noreplace) %attr(0644,root,root) %{_sysconfdir}/system-services-ui/system-services-ui.yaml

%changelog
* Sat Dec 27 2025 Andrei Ivan <contact@andreiivan.com> - 0.1.0-1
- Initial Fedora package
