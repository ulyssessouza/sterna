#!/usr/bin/env bash

set -e
set -xv

./sql_dump.sh ${BLUE_HOST} ${BLUE_PORT} dumped-data.sql ${DB}
./sql_execute.sh ${GREEN_HOST} ${GREEN_PORT} create_db.sql
./sql_execute.sh ${GREEN_HOST} ${GREEN_PORT} dumped-data.sql ${DB}
./sql_execute.sh ${GREEN_HOST} ${GREEN_PORT} alter_table.sql ${DB}

#rm -rf dumped-data.sql
