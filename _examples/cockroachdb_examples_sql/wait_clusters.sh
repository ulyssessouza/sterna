#!/usr/bin/env bash

set -e
set -xv

while [ "$(cockroach node --insecure --host ${GREEN_HOST} -p ${GREEN_PORT} ls --format csv | wc -l)" != "5" ]
do
    sleep 1
done

sleep 10