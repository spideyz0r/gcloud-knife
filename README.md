# gcloud-knife [![CI](https://github.com/spideyz0r/gcloud-knife/workflows/gotester/badge.svg)][![CI](https://github.com/spideyz0r/gcloud-knife/workflows/goreleaser/badge.svg)][![CI](https://github.com/spideyz0r/gcloud-knife/workflows/rpm-builder/badge.svg)]
gcloud-knife is a tool to run commands on multiple GCP instances in concurrently.

It can also be used to initially list the instances to ensure that the filter is working correctly before executing commands on the VMs.


## Install

### RPM
```
dnf copr enable brandfbb/gcloud-knife
dnf install gcloud-knife
```


### From source
```
go build -v -o  gcloud-knife
```

## Usage
```
# gcloud-knife
Usage: gcloud-knife [-h] [-c value] [-f value] [-p value] [-t value] [-u value] [parameters ...]
 -c, --command=value
                   command to run, if empty will print the list of hosts
 -f, --filter=value
                   filter
 -h, --help        display this help
 -p, --project=value
                   project
 -t, --port=value  optional port, default: 22
 -u, --user=value  optional user, otherwise will read from local configuration
```

## Examples
Execute a command in multiple VMs filtering by their name
```
gcloud-knife -p my-gcp-project -f name:*webservers.mydomain.com -c "df -H | grep nfs"
```

List VMs that match a given filter
```
gcloud-knife -p my-gcp-project -f region:someregion
```
