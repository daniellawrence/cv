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


for SERVICE_DOCKERFILE in  $(find frontend -name Dockerfile);do
    SERVICE=$(dirname $SERVICE_DOCKERFILE)
    [[ -n "${SERVICE_FILTER}" && "${SERVICE}" != *"${SERVICE_FILTER}"* ]] && continue
    CONTENT_SHA1=$(frontend_content_sha ${SERVICE})
    IMAGE_NAME=${IMAGE_PREFIX}/${SERVICE}:latest
    IMAGE_FILE="dist/${SERVICE}.${CONTENT_SHA1}.tar.gz"

    if [[ -f ${IMAGE_FILE} ]];then
        echo "image: ${SERVICE} - ${IMAGE_FILE}"
    else
        rm "dist/${SERVICE}.*.tar.gz"

        docker build -t ${IMAGE_NAME} -f ${SERVICE_DOCKERFILE} .
        docker save ${IMAGE_NAME} | gzip > ${IMAGE_FILE}
    fi

    if [[ "${TARGET}" == "dev" ]];then
        ./bin/k3d image import --cluster ${KUBE_CLUSER_NAME} ${IMAGE_FILE}
        docker exec -it k3d-k0-server-0 crictl images |grep $SERVICE
    fi

done

for SERVICE_DOCKERFILE in $(find backend -name Dockerfile);do
    SERVICE=$(basename $(dirname $SERVICE_DOCKERFILE))
    [[ -n "${SERVICE_FILTER}" && "${SERVICE}" != *"${SERVICE_FILTER}"* ]] && continue
    CONTENT_SHA1=$(backend_content_sha ${SERVICE})
    IMAGE_NAME=${IMAGE_PREFIX}/${SERVICE}:latest
    IMAGE_FILE="dist/${SERVICE}.${CONTENT_SHA1}.tar.gz"

    if [[ -f ${IMAGE_FILE} ]];then
        echo "image: ${SERVICE} - ${IMAGE_FILE}"
    else
        rm "dist/${SERVICE}.*.tar.gz"
        docker build -t ${IMAGE_NAME} -f ${SERVICE_DOCKERFILE} .
        docker save ${IMAGE_NAME} | gzip > ${IMAGE_FILE}
    fi

    if [[ "${TARGET}" == "dev" ]];then
        ./bin/k3d image import --cluster ${KUBE_CLUSER_NAME} ${IMAGE_FILE}
        docker exec -it k3d-k0-server-0 crictl images |grep $SERVICE
    fi

done

if [[ "${TARGET}" == "prod" ]];then
    rsync -e ssh -avz dist/* root@dansysadm.com:/dist/
    for IMAGE in $(cd ./dist && ls );do
        ssh root@dansysadm.com "cat /dist/${IMAGE} | gunzip | k3s ctr images import -"
    done
    ssh root@dansysadm.com "k3s ctr images ls |grep dansysadm.com/images/"
fi