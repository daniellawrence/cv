#!/usr/bin/env bash
source $(dirname $0)/lib.sh
cd $GIT_ROOT

TARGET="dev"
SERVICE_FILTER=""

for arg in "$@"; do
    case "$arg" in
        prod) TARGET="prod" ;;
        --service) shift; SERVICE_FILTER="$1" ;;
        --service=*) SERVICE_FILTER="${arg#--service=}" ;;
    esac
    shift 2>/dev/null || true
done

if [[ "${TARGET}" == "prod" ]]; then
    IMAGE_PREFIX="localhost:5001"
    if ! curl -sf http://localhost:5001/v2/ > /dev/null 2>&1; then
        echo "ERROR: production registry not reachable at localhost:5001"
        echo "Open the SSH tunnel first:"
        echo "  ssh -NL 7443:localhost:6443 -L 5001:localhost:5001 root@dansysadm.com"
        exit 1
    fi
fi

BUILD_TIMESTAMP=$(date -u +%Y-%m-%dT%H:%M:%SZ)

for SERVICE_DOCKERFILE in  $(find frontend -name Dockerfile);do
    SERVICE=$(dirname $SERVICE_DOCKERFILE)
    [[ -n "${SERVICE_FILTER}" && "${SERVICE}" != *"${SERVICE_FILTER}"* ]] && continue
    CONTENT_SHA1=$(frontend_content_sha ${SERVICE})
    IMAGE_NAME=${IMAGE_PREFIX}/${SERVICE}:latest
    
    docker build -t ${IMAGE_NAME} -f ${SERVICE_DOCKERFILE} \
        --build-arg CONTENT_SHA1="${CONTENT_SHA1}" \
        .
    docker push ${IMAGE_NAME}

done

for SERVICE_DOCKERFILE in $(find backend -name Dockerfile);do
    SERVICE=$(basename $(dirname $SERVICE_DOCKERFILE))
    [[ -n "${SERVICE_FILTER}" && "${SERVICE}" != *"${SERVICE_FILTER}"* ]] && continue
    CONTENT_SHA1=$(backend_content_sha ${SERVICE})
    IMAGE_NAME=${IMAGE_PREFIX}/${SERVICE}:latest
    
    docker build -t ${IMAGE_NAME} -f ${SERVICE_DOCKERFILE} \
        --build-arg CONTENT_SHA1="${CONTENT_SHA1}" \
        .
    docker push ${IMAGE_NAME}

done