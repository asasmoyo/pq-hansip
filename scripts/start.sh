#!/usr/bin/env bash

source './scripts/_common.sh'

for host in $hosts; do
    current="${prefix}-${host}"

    if docker ps --format '{{ .Names }}' | grep --quiet "${current}"; then
        echo "${host} already running!"
        continue
    fi

    echo "running '${host}' database..."
    docker run --name="${current}" -e POSTGRES_DB="${db_name}" -e POSTGRES_USER="${db_user}" -e POSTGRES_PASSWORD="${db_pass}" --publish-all --rm --detach postgres:11 &> /dev/null
done

rm -f .env
for host in $hosts; do
    key="DB_${host^^}_URL"
    port=$(docker inspect "${prefix}-${host}" --format '{{ (index .NetworkSettings.Ports "5432/tcp" 0).HostPort }}')
    val="postgres://${db_user}:${db_pass}@localhost:${port}/${prefix}?sslmode=disable"
    echo "export ${key}=${val}" >> .env
done
