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

#go vet -v ./...
#go build -v ./...
go test -v aggregate/* -cover
go test -v logger/* -cover
go test -v tlscerts/* -cover
go test -v tokenbucket/* -cover
go test -v units/* -cover

exit $?
