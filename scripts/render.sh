#!/usr/bin/env bash
source $(dirname $0)/lib.sh
cd $GIT_ROOT


for ENVIRONMENT in "dev" "production"; do

    ./bin/helmfile  \
        --environment ${ENVIRONMENT} \
        --helm-binary ${PWD}/bin/helm \
        --file infra/deployments/helmfile.yaml \
        --selector tier=tracing \
        template \
        --output-dir-template "artifacts/${ENVIRONMENT}/{{ .Release.Name }}/"

    ./bin/helmfile  \
        --environment ${ENVIRONMENT} \
        --helm-binary ${PWD}/bin/helm \
        --file infra/deployments/helmfile.yaml \
        --selector tier=app \
        template \
        --output-dir-template "artifacts/${ENVIRONMENT}/{{ .Release.Name }}/"

    if [[ "${ENVIRONMENT}" == "dev" ]];then

        ./bin/helmfile \
            --environment ${ENVIRONMENT} \
            --helm-binary ${PWD}/bin/helm \
            --file infra/deployments/helmfile.yaml \
            --selector tier=deploy \
            template \
            --output-dir-template "artifacts/${ENVIRONMENT}/{{ .Release.Name }}/"
        fi

done