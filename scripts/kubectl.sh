#!/usr/bin/env bash
source $(dirname $0)/lib.sh
cd $GIT_ROOT
bin/kubectl --context k3d-k0 ${@}