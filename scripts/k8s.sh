#!/usr/bin/env bash
source $(dirname $0)/lib.sh
cd $GIT_ROOT
./bin/k3d cluster create k0 --config infra/k8s/k3d.development.yaml
./bin/k3d kubeconfig write k0 --output=$KUBECONFIG