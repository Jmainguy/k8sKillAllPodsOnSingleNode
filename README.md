# k8sKillAllPodsOnSingleNode
Force kill all pods on a single k8s node.

## Build
```/bin/bash
export GO111MODULE=on
go mod init
go get k8s.io/client-go@v12.0.0
go build
```

## Usage
Login to your kubernetes cluster, then run
```/bin/bash
./k8sKillAllPodsOnSingleNode --nodeName ocp-app-01a.example.com
```


### OC equivalent
You can bypass this tool and just use an oc command like
```/bin/bash
oc adm drain ocp-app-01a.example.com --force --grace-period=0 --delete-local-data --ignore-daemonsets
```
