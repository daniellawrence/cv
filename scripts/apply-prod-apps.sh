#!/usr/bin/env bash
# Applies ArgoCD prod Application manifests to the ArgoCD namespace.
# Called by the helmfile postsync hook after ArgoCD is deployed.
set -euo pipefail
source "$(dirname "$0")/lib.sh"
cd "$GIT_ROOT"

find infra/argocd/prod -name '*.yaml' | while read -r f; do
    ./bin/kubectl apply -f "$f"
done
