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

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
export GOPATH="$script_dir"
export GOBIN="${GOPATH}/bin"
#go env

printf "Downloading and installing packages and dependencies...\n"

go get -v . 

printf "Compiling packages and dependencies...\n"
go build -v ./...

go test -v ./...

exit $?
