#/bin/sh -xe

POD_PORT=$(kubectl get gs -o json | jq -r '.items[0].status.ports[0] | select(.name == "default").port')
POD_NAME=$(kubectl get pod -o json | jq -r '.items[0].metadata.name')

kubectl port-forward $POD_NAME 7654:$POD_PORT