#!/bin/bash

set -e

echo "Starting Minikube"
minikube start --driver=docker

echo "Enable Ingress"
minikube addons enable ingress

echo "Deploying Helm Chart"
./install-helm.sh

echo "Getting service URLs..."
minikube service list
