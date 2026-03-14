#!/usr/bin/env bash
# Refreshes the ArgoCD cluster secret for the production cluster.
# Run this after the production k3d cluster is recreated or certs rotate.
#   ./scripts/sync-prod-cluster-creds.sh
set -euo pipefail
source "$(dirname "$0")/lib.sh"
cd "$GIT_ROOT"

PROD_SERVER="https://k8s.dansysadm.com:6443"

# --- 1. Validate kubeconfig is reachable ---------------------------------
echo "==> Checking production cluster connectivity..."
if ! ./bin/kubectl --kubeconfig "${PROD_KUBECONFIG}" get nodes > /dev/null 2>&1; then
    echo "ERROR: cannot reach production cluster at ${PROD_SERVER}"
    echo "  Run: cd infra/ansible && make k3d"
    exit 1
fi
echo "    OK"

# --- 2. Extract creds from kubeconfig ------------------------------------
echo "==> Extracting credentials from ${PROD_KUBECONFIG}..."
CA_DATA=$(./bin/kubectl --kubeconfig "${PROD_KUBECONFIG}" config view --raw \
    -o jsonpath='{.clusters[0].cluster.certificate-authority-data}')
CERT_DATA=$(./bin/kubectl --kubeconfig "${PROD_KUBECONFIG}" config view --raw \
    -o jsonpath='{.users[0].user.client-certificate-data}')
KEY_DATA=$(./bin/kubectl --kubeconfig "${PROD_KUBECONFIG}" config view --raw \
    -o jsonpath='{.users[0].user.client-key-data}')

if [[ -z "${CA_DATA}" || -z "${CERT_DATA}" || -z "${KEY_DATA}" ]]; then
    echo "ERROR: kubeconfig is missing credentials (CA, cert, or key is empty)"
    echo "  Run: cd infra/ansible && make k3d"
    exit 1
fi
echo "    OK"

# --- 3. Patch the ArgoCD cluster secret ----------------------------------
echo "==> Patching cluster-production secret in ArgoCD..."
CONFIG=$(printf '{"tlsClientConfig":{"insecure":false,"caData":"%s","certData":"%s","keyData":"%s"}}' \
    "${CA_DATA}" "${CERT_DATA}" "${KEY_DATA}")

./bin/kubectl patch secret cluster-production -n argocd \
    --type=merge \
    -p "{\"stringData\":{\"server\":\"${PROD_SERVER}\",\"config\":$(printf '%s' "${CONFIG}" | python3 -c 'import json,sys; print(json.dumps(sys.stdin.read()))')}}"
echo "    OK"

# --- 4. Validate ArgoCD can now reach the cluster ------------------------
echo "==> Waiting for ArgoCD to reconnect..."
sleep 5

for i in $(seq 1 12); do
    STATUS=$(./bin/kubectl get secret cluster-production -n argocd \
        -o jsonpath='{.metadata.annotations.argocd\.argoproj\.io/cluster-connection-state}' 2>/dev/null || true)

    CLUSTER_READY=$(./bin/kubectl get applications -n argocd \
        -o jsonpath='{.items[*].status.conditions[?(@.type=="SyncError")].message}' 2>/dev/null | \
        grep -c "certificate signed by unknown authority" || true)

    if [[ "${CLUSTER_READY}" == "0" ]]; then
        echo "    OK — ArgoCD cluster connection restored"
        exit 0
    fi
    echo "    Waiting... (${i}/12)"
    sleep 5
done

echo "WARNING: could not confirm ArgoCD reconnection — check the ArgoCD UI for sync errors"
exit 1
