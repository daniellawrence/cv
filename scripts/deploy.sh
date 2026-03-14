#!/usr/bin/env bash
source $(dirname $0)/lib.sh
cd $GIT_ROOT

export ENVIRONMENT="dev"
if [[ "${1}" == "prod" ]];then
    export ENVIRONMENT="production"
    export KUBECONFIG=${GIT_ROOT}/infra/k8s/kubeconfig.production.yaml
    if ! ./bin/kubectl get nodes > /dev/null 2>&1; then
        echo "ERROR: cannot reach production cluster at https://k8s.dansysadm.com:6443"
        exit 1
    fi
    ./scripts/make_images.sh prod
fi

./bin/helmfile  \
    --environment ${ENVIRONMENT} \
    --helm-binary ${PWD}/bin/helm \
    --file infra/deployments/helmfile.yaml \
    template \
    --output-dir-template "artifacts/${ENVIRONMENT}/{{ .Release.Name }}/" 

if [[ "${ENVIRONMENT}" == "dev" ]];then
    ./bin/helmfile \
        --environment ${ENVIRONMENT} \
        --helm-binary ${PWD}/bin/helm \
        --file infra/deployments/helmfile.yaml \
        --selector tier=infra \
        sync
    ./bin/helmfile \
        --environment ${ENVIRONMENT} \
        --helm-binary ${PWD}/bin/helm \
        --file infra/deployments/helmfile.yaml \
        --selector tier=deploy \
        sync

fi

./bin/kubectl rollout restart deployment -n education
./bin/kubectl rollout restart deployment -n experience
./bin/kubectl rollout restart deployment -n identity
./bin/kubectl rollout restart deployment -n interest
./bin/kubectl rollout restart deployment -n qrcode
./bin/kubectl rollout restart deployment -n frontend