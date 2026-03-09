#!/usr/bin/env bash
source $(dirname $0)/lib.sh
cd $GIT_ROOT

IMAGE_PREFIX="dansysadm.com/images"

for SERVICE_DOCKERFILE in $(find backend -name Dockerfile);do
    SERVICE=$(basename $(dirname $SERVICE_DOCKERFILE))
    IMAGE_NAME=${IMAGE_PREFIX}/${SERVICE}:latest
    docker build -t ${IMAGE_NAME} -f ${SERVICE_DOCKERFILE} .
    ./bin/k3d image import --cluster ${KUBE_CLUSER_NAME} ${IMAGE_NAME}
    echo $IMAGE_NAME
done

