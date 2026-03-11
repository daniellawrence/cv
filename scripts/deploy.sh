#!/usr/bin/env bash
source $(dirname $0)/lib.sh
cd $GIT_ROOT

export ENVIRONMENT="dev"
if [[ "${1}" == "prod" ]];then
    export ENVIRONMENT="production"
    export KUBECONFIG=${GIT_ROOT}/infra/k8s/k3s.yaml
fi

./bin/helmfile  \
    --environment ${ENVIRONMENT} \
    --helm-binary ${PWD}/bin/helm \
    --file infra/deployments/helmfile.yaml \
    template \
    --output-dir-template "artifacts/${ENVIRONMENT}/{{ .Release.Name }}/" 

./bin/helmfile  \
    --environment ${ENVIRONMENT} \
    --helm-binary ${PWD}/bin/helm \
    --file infra/deployments/helmfile.yaml \
    sync

./bin/kubectl rollout restart deployment -n education
./bin/kubectl rollout restart deployment -n experience
./bin/kubectl rollout restart deployment -n identity
./bin/kubectl rollout restart deployment -n interest
./bin/kubectl rollout restart deployment -n qrcode
./bin/kubectl rollout restart deployment -n frontend