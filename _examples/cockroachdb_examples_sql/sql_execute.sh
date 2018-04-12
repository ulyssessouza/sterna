#!/usr/bin/env bash

set -xv

HOST="${1}"
PORT="${2}"
FILE="${3}"
DB="${4}"

if [ "${DB}" != "" ]; then
    DB_OPTION="--database ${DB}"
fi

cockroach sql --insecure --host ${HOST} -p ${PORT} ${DB_OPTION} < ${FILE}