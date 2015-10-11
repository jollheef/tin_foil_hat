#!/bin/sh
START_TIME=`date +%s`

COMMIT_ID=`git --no-pager log --format="%H" -n 1`
BUILD_DATE=`date -u +%d.%m.%Y`
BUILD_TIME=`date -u +%H:%M:%S`

LDFLAGS="-X main.COMMIT_ID ${COMMIT_ID}"
LDFLAGS+=" -X main.BUILD_DATE ${BUILD_DATE}"
LDFLAGS+=" -X main.BUILD_TIME ${BUILD_TIME}"

export GOPATH=$(realpath ./)

REPO=github.com/jollheef/tin_foil_hat/

go build -ldflags "${LDFLAGS}" -o bin/tinfoilhat src/${REPO}/tinfoilhat.go
go build -ldflags "${LDFLAGS}" -o bin/tfhctl src/${REPO}/tfhctl.go

END_TIME=`date +%s`
RUN_TIME=$((END_TIME-START_TIME))
echo "Build done in ${RUN_TIME}s"
