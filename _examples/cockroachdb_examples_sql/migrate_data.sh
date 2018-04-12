#!/usr/bin/env bash

set -e
set -xv

./sql_dump.sh ${BLUE_HOST} ${BLUE_PORT} dump.sql ${DB}
./sql_execute.sh ${GREEN_HOST} ${GREEN_PORT} create_db.sql
./sql_execute.sh ${GREEN_HOST} ${GREEN_PORT} dump.sql ${DB}
