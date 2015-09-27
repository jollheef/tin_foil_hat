#!/bin/sh

START_TIME=`date +%s`

PKGS=`find src/tinfoilhat -mindepth 1 -maxdepth 1 -type d | sed 's/src\///'`

for PKG in ${PKGS}; do
    GOPATH=$(realpath ./) go test ${PKG}
done

END_TIME=`date +%s`

RUN_TIME=$((END_TIME-START_TIME))

echo "All done in ${RUN_TIME}s"
