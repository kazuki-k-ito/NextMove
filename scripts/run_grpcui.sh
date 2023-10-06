#!/bin/sh -xe

TOP_DIR=$(git rev-parse --show-toplevel)

PORT=$1

#portの解放
POD_NAME=$(kubectl get po --no-headers -o custom-columns=":metadata.name" | tr -d '\n')
kubectl port-forward $POD_NAME $PORT:$PORT &

sleep 3

grpcui -import-path $TOP_DIR/schema -proto game.proto -plaintext localhost:$PORT
