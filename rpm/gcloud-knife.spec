%global go_version 1.18.10
%global go_release go1.18.10

Name:           gcloud-knife
Version:        0.3
Release:        1%{?dist}
Summary:        This is a test
License:        GPLv3
URL:            https://github.com/spideyz0r/gcloud-knife
Source0:        %{url}/archive/refs/tags/v%{version}.tar.gz

BuildRequires:  golang >= %{go_version}
BuildRequires:  git

%description
gcloud-knife is a tool to run commands on multiple GCP instances in parallel.
It can also be used to initially list the instances to ensure that the filter is working correctly before executing commands on the VMs.

%global debug_package %{nil}

%prep
%autosetup -n %{name}-%{version}

%build
go build -v -o %{name} -ldflags=-linkmode=external

%check
go test

%install
install -Dpm 0755 %{name} %{buildroot}%{_bindir}/%{name}


%files
%{_bindir}/gcloud-knife

%license LICENSE

%changelog
* Fri Mar 24 2023 spideyz0r <47341410+spideyz0r@users.noreply.github.com> 0.3-1
- Initial build

* Fri Mar 24 2023 spideyz0r <47341410+spideyz0r@users.noreply.github.com> 0.1-1
- Initial build
