#
# spec file for package sql_exporter
#
# Copyright (c) 2018 credativ GmbH
#

Name:			sql_exporter
Version:		0.2.0
Release:		1%{?dist}
License:		MIT
Summary:		Exporter for SQL metrics
Url:			https://%{provider_prefix}
Group:			System/Monitoring
Source0:		%{name}-%{version}.tar.gz
Source1:		prometheus-sql-exporter.default
Source2:		prometheus-sql-exporter.service
Patch1:			synchronous-jobs
BuildRequires:	golang
BuildRequires:  go
#BuildRequires:  golang-github-prometheus-promu
BuildRequires:  systemd
BuildRequires:  help2man
%{?systemd_requires}
BuildRoot:      %{_tmppath}/%{name}-%{version}-build

%description
Prometheus exporter for SQL metrics in Go with pluggable metric collectors.

%prep
%setup -q -n %{name}-%{version}
%patch1 -p1

%build
GOPATH=%{_builddir}/go go get -v github.com/credativ/sql_exporter
GOPATH=%{_builddir}/go go build -x -v github.com/credativ/sql_exporter

%install
install -m 0755 -d %{buildroot}%{_sbindir}
install -m 0755 -D -s sql_exporter %{buildroot}%{_bindir}/prometheus-sql-exporter
install -m 0644 -D %{_sourcedir}/prometheus-sql-exporter.default %{buildroot}%{_sysconfdir}/default/prometheus-sql-exporter
install -m 0644 -D %{_sourcedir}/prometheus-sql-exporter.service %{buildroot}%{_unitdir}/prometheus-sql-exporter.service

# manpage
mkdir -p %{buildroot}%{_mandir}/man1
PATH=%{buildroot}%{_bindir}:$PATH help2man --no-discard-stderr --no-info --name "SQL Exporter for Prometheus" --version-string=%{version} prometheus-sql_exporter > prometheus-sql-exporter.1
install -m 0644 -D prometheus-sql-exporter.1 %{buildroot}%{_mandir}/man1/prometheus-sql-exporter.1

%files
%defattr(-,root,root)
%doc LICENSE
%{_bindir}/prometheus-sql-exporter
%{_unitdir}/prometheus-sql-exporter.service
%config %{_sysconfdir}/default/prometheus-sql-exporter
%{_mandir}/man1/prometheus-sql-exporter.1.gz

%changelog
