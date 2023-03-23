# gcloud-knife
gcloud-knife is a tool to run commands on multiple GCP instances in parallel.

It can also be used to initially list the instances to ensure that the filter is working correctly before executing commands on the VMs.

```
# gcloud-knife
Usage: gcloud-knife [-h] [-c value] [-f value] [-p value] [-u value] [parameters ...]
 -c, --command=value
                   command to be run externally
 -f, --filter=value
                   target filter
 -h, --help        display this help
 -p, --project=value
                   GCP project
 -u, --user=value  [user]
 ```
