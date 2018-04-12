#!/usr/bin/env bash

HOST="192.168.99.100"
PORT="30251"
DB="test"

pushd _examples/cockroachdb_examples_sql
    ./sql_execute.sh ${HOST} ${PORT} create_db.sql
    ./sql_execute.sh ${HOST} ${PORT} create_table.sql ${DB}
    ./sql_execute.sh ${HOST} ${PORT} insert.sql ${DB}
popd