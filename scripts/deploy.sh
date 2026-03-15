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

# Step 1: Render charts to verify they produce valid output
./bin/helmfile  \
    --environment ${ENVIRONMENT} \
    --helm-binary ${PWD}/bin/helm \
    --file infra/deployments/helmfile.yaml \
    template \
    --output-dir-template "artifacts/${ENVIRONMENT}/{{ .Release.Name }}/" 

echo "✓ Chart rendering completed successfully"

# Step 2: Run chart tests before syncing (always, for all environments)
echo ""
echo "Running chart tests..."

# Check if any releases are already deployed using helmfile ls
DEPLOYED_RELEASES=$(./bin/helmfile --environment ${ENVIRONMENT} \
    --helm-binary ${PWD}/bin/helm \
    --file infra/deployments/helmfile.yaml \
    ls 2>/dev/null | tail -n +2 || true)

if [ -n "${DEPLOYED_RELEASES}" ]; then
    echo "Found deployed releases, running helmfile test..."
    
    # Run helmfile test which will test all deployed releases
    if ./bin/helmfile --environment ${ENVIRONMENT} \
        --helm-binary ${PWD}/bin/helm \
        --file infra/deployments/helmfile.yaml \
        test 2>&1; then
        echo ""
        echo "✓ All chart tests passed successfully"
    else
        echo ""
        echo "✗ Chart tests failed - deployment aborted"
        exit 1
    fi
else
    echo "No deployed releases found, skipping automatic testing."
    echo "Run './scripts/deploy.sh test' to manually test specific charts."
fi

echo ""

# Step 3: Sync infra tier deployments
./bin/helmfile \
    --environment ${ENVIRONMENT} \
    --helm-binary ${PWD}/bin/helm \
    --file infra/deployments/helmfile.yaml \
    --selector tier=infra \
    sync
    
if [[ "${ENVIRONMENT}" == "dev" ]];then

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

# Function to run chart tests for cv-app
run-chart-tests() {
    local release_name="${1:-cv-app}"
    local namespace="${2:-education}"
    
    echo "Running chart tests for ${release_name} in namespace ${namespace}..."
    
    if ./bin/helm test "${release_name}" -n "${namespace}" 2>&1; then
        echo ""
        echo "✓ Chart tests passed successfully"
        return 0
    else
        echo ""
        echo "✗ Chart tests failed"
        return 1
    fi
}

# If called with 'test' argument, run chart tests instead of deploying
if [[ "${1}" == "test" ]]; then
    release_name="${2:-cv-app}"
    namespace="${3:-education}"
    run-chart-tests "${release_name}" "${namespace}"
    exit $?
fi