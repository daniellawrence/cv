GIT_ROOT=$(git rev-parse --show-toplevel)
BIN_DIR="${GIT_ROOT}/bin"
export KUBECONFIG=${GIT_ROOT}/infra/k8s/kubeconfig.dev.yaml
export PROD_KUBECONFIG=${GIT_ROOT}/infra/k8s/kubeconfig.production.yaml
export KUBE_CLUSER_NAME=k0
export IMAGE_PREFIX="localhost:5000"


function backend_content_sha() {
    SERVICE_NAME=$1
    BACKEND_SHA1=$(
        (
            find backend/${SERVICE_NAME}/ -type f -exec cat {} \;
            find backend/common -type f -exec cat {} \;
            find gen/go -type f -exec cat {} \; 
        ) | sha1sum | cut -d' ' -f1)
    echo $BACKEND_SHA1
}

function frontend_content_sha() {
    FRONTEND_SHA1=$(
        (
            find frontend/ -type f -exec cat {} \;
            find gen/ts -type f -exec cat {} \; 
        ) | sha1sum | cut -d' ' -f1)
    echo $FRONTEND_SHA1
}