#!/bin/bash

set -e
set -o errtrace

function err_handler() {
    local frame=0
    while caller $frame; do
        ((frame++));
    done
    echo "$*"
    exit 1
}

trap 'err_handler' SIGINT ERR

go build -v ./...
go test -v ./... -cover

exit $?
