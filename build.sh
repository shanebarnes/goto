#!/bin/bash

set -e
set -o errtrace

function err_handler() {
    local frame=0
    while caller $frame; do
        ((frame++));
    done
    printf "%s\n" "$*"
    exit 1
}

trap 'err_handler' SIGINT ERR

gofmt_diff=$(gofmt -d .)
if [[ ! -z "$gofmt_diff" ]]; then
    printf "%s\n" "$gofmt_diff"
    printf "To continue building, please fix Go formatting by running 'gofmt -w .'\n"
    exit 1
fi

go vet -v ./...
go build -v ./...
go test -v ./... -cover

exit $?
