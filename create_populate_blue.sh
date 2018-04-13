#!/usr/bin/env bash

HOST="192.168.99.100"
PORT="30251"
DB="test"

pushd _examples/cockroachdb_examples_sql
    # Create a database with name "test"
    ./sql_execute.sh ${HOST} ${PORT} create_db.sql

    # Create a table with name "test_tb" in the created database
    ./sql_execute.sh ${HOST} ${PORT} create_table.sql ${DB}

    # Insert one row containing "Ulysses" in the column "name"
    ./sql_execute.sh ${HOST} ${PORT} insert.sql ${DB}
popd