#!/bin/sh
START_TIME=`date +%s`

RUNDIR=$(pwd)

cd $(dirname ${0})
COMMIT_ID=`git --no-pager log --format="%H" -n 1`

cd ${RUNDIR}

BUILD_DATE=`date -u +%d.%m.%Y`
BUILD_TIME=`date -u +%H:%M:%S`

LDFLAGS="-X main.COMMIT_ID ${COMMIT_ID}"
LDFLAGS+=" -X main.BUILD_DATE ${BUILD_DATE}"
LDFLAGS+=" -X main.BUILD_TIME ${BUILD_TIME}"

export GOPATH=$(realpath ./)

mkdir -p bin

go build -ldflags "${LDFLAGS}" -o bin/tin_foil_hat github.com/jollheef/tin_foil_hat
go build -ldflags "${LDFLAGS}" -o bin/tfhctl github.com/jollheef/tin_foil_hat/cli/tfhctl

END_TIME=`date +%s`
RUN_TIME=$((END_TIME-START_TIME))
echo "Build done in ${RUN_TIME}s"
