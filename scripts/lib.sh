GIT_ROOT=$(git rev-parse --show-toplevel)
BIN_DIR="${GIT_ROOT}/bin"
export KUBECONFIG=${GIT_ROOT}/infra/k8s/kubeconfig.yaml
export PROD_KUBECONFIG=${GIT_ROOT}/infra/k8s/k3s.yaml
export KUBE_CLUSER_NAME=k0
IMAGE_PREFIX="dansysadm.com/images"
