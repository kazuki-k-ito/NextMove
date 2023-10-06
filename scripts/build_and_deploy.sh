#!/bin/sh -xe

pushd $(git rev-parse --show-toplevel)

TAG=$(openssl rand -hex 16)

# Build the docker image
pushd server
docker build . -t localimage/grpc-server:$TAG
minikube image load localimage/grpc-server:$TAG
popd

# Kubernetes apply
CURRENT_CONTEXT=$(kubectl config current-context)
if [ "$CURRENT_CONTEXT" != "minikube" ]; then
  echo "Please set your kubectl context to minikube"
  exit 1
fi

cat agones/fleet.yaml | sed -e "s/:latest/:${TAG}/g" | kubectl apply -f -

popd
