#!/usr/bin/env bash
source $(dirname $0)/lib.sh
cd $GIT_ROOT

./bin/helmfile  \
    --helm-binary ${PWD}/bin/helm \
    --file infra/deployments/helmfile.yaml \
    sync