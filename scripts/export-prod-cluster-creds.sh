#!/usr/bin/env bash
# Exports production cluster credentials as environment variables for helmfile.
# Run this with `source` before deploying ArgoCD:
#   source scripts/export-prod-cluster-creds.sh
set -euo pipefail
source "$(dirname "$0")/lib.sh"
cd "$GIT_ROOT"

if [[ ! -f "${PROD_KUBECONFIG}" ]]; then
    echo "ERROR: production kubeconfig not found at ${PROD_KUBECONFIG}"
    exit 1
fi

if ! ./bin/kubectl --kubeconfig "${PROD_KUBECONFIG}" cluster-info > /dev/null 2>&1; then
    echo "ERROR: production cluster not reachable at https://k8s.dansysadm.com:6443"
    exit 1
fi

export ARGOCD_PROD_CA_DATA
export ARGOCD_PROD_CERT_DATA
export ARGOCD_PROD_KEY_DATA

ARGOCD_PROD_CA_DATA=$(KUBECONFIG="${PROD_KUBECONFIG}" ./bin/kubectl config view --raw \
    -o jsonpath='{.clusters[0].cluster.certificate-authority-data}')
ARGOCD_PROD_CERT_DATA=$(KUBECONFIG="${PROD_KUBECONFIG}" ./bin/kubectl config view --raw \
    -o jsonpath='{.users[0].user.client-certificate-data}')
ARGOCD_PROD_KEY_DATA=$(KUBECONFIG="${PROD_KUBECONFIG}" ./bin/kubectl config view --raw \
    -o jsonpath='{.users[0].user.client-key-data}')

echo "Production cluster credentials exported. Run deploy.sh to apply."
