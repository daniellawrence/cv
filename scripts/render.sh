#!/usr/bin/env bash
source $(dirname $0)/lib.sh
cd $GIT_ROOT


for ENVIRONMENT in "dev" "production"; do

    ./bin/helmfile  \
        --environment ${ENVIRONMENT} \
        --helm-binary ${PWD}/bin/helm \
        --file infra/deployments/helmfile.yaml \
        template \
        --output-dir-template "artifacts/${ENVIRONMENT}/{{ .Release.Name }}/" 

done