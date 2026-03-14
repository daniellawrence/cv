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
    IMAGE_PREFIX="k8s.dansysadm.com:5001"
    probe_result=$(docker manifest inspect "${IMAGE_PREFIX}/probe:test" 2>&1 || true)
    if echo "${probe_result}" | grep -qiE "server gave HTTP response|connection refused|no such host|unable to resolve"; then
        echo "ERROR: docker cannot reach registry at ${IMAGE_PREFIX}"
        echo "  Ensure ${IMAGE_PREFIX} is listed in /etc/docker/daemon.json insecure-registries"
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