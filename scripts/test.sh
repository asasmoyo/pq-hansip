#!/usr/bin/env bash

source './scripts/_common.sh'
trap './scripts/stop.sh' EXIT
./scripts/start.sh
sleep 3

source .env
go test -v -cover -race .
