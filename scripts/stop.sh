#!/usr/bin/env bash

source './scripts/_common.sh'

for host in $hosts; do
    current="${prefix}-${host}"

    if ! $(docker ps --format '{{ .Names }}' | grep --quiet $current); then
        echo "${host} is not running!"
        continue
    fi

    echo "stopping '${host}' database..."
    docker kill "${current}" &> /dev/null
done
