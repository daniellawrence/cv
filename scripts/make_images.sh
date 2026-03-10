#!/usr/bin/env bash
source $(dirname $0)/lib.sh
cd $GIT_ROOT

TARGET="dev"
if [[ "${1}" == "prod" ]];then
    TARGET="prod"
fi


for SERVICE_DOCKERFILE in  $(find frontend -name Dockerfile);do
    SERVICE=$(dirname $SERVICE_DOCKERFILE)
    CONTENT_SHA1=$(find ${SERVICE}/ -type f -exec cat {} \; | sha1sum | cut -d' ' -f1)
    IMAGE_NAME=${IMAGE_PREFIX}/${SERVICE}:latest
    IMAGE_FILE="dist/${SERVICE}.${CONTENT_SHA1}.tar.gz"

    if [[ -f ${IMAGE_FILE} ]];then
        echo "image: ${SERVICE} - ${IMAGE_FILE}"
        continue
    fi
    rm "dist/${SERVICE}.*.tar.gz"

    docker build -t ${IMAGE_NAME} -f ${SERVICE_DOCKERFILE} .
    docker save ${IMAGE_NAME} | gzip > ${IMAGE_FILE}

    if [[ "${TARGET}" == "dev" ]];then
        ./bin/k3d image import --cluster ${KUBE_CLUSER_NAME} ${IMAGE_NAME}
    fi

done

for SERVICE_DOCKERFILE in $(find backend -name Dockerfile);do
    SERVICE=$(basename $(dirname $SERVICE_DOCKERFILE))
    IMAGE_NAME=${IMAGE_PREFIX}/${SERVICE}:latest
    CONTENT_SHA1=$(find backend/${SERVICE}/ -type f -exec cat {} \; | sha1sum | cut -d' ' -f1)
    IMAGE_NAME=${IMAGE_PREFIX}/${SERVICE}:latest
    IMAGE_FILE="dist/${SERVICE}.${CONTENT_SHA1}.tar.gz"

    if [[ -f ${IMAGE_FILE} ]];then
        echo "image: ${SERVICE} - ${IMAGE_FILE}"
        continue
    fi
    rm "dist/${SERVICE}.*.tar.gz"

    docker build -t ${IMAGE_NAME} -f ${SERVICE_DOCKERFILE} .
    docker save ${IMAGE_NAME} | gzip > ${IMAGE_FILE}

    if [[ "${TARGET}" == "dev" ]];then
        ./bin/k3d image import --cluster ${KUBE_CLUSER_NAME} ${IMAGE_NAME}
    fi

done

if [[ "${TARGET}" == "prod" ]];then
    rsync -e ssh -avz dist/* root@dansysadm.com:/dist/
    for IMAGE in $(cd ./dist && ls );do
        ssh root@dansysadm.com "cat /dist/${IMAGE} | gunzip | k3s ctr images import -"
    done
    ssh root@dansysadm.com "k3s ctr images ls |grep dansysadm.com/images/"
fi