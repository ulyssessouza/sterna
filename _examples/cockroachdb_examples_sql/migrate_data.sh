#!/usr/bin/env bash

set -e
set -xv

# Dump data from BLUE to a file
./sql_dump.sh ${BLUE_HOST} ${BLUE_PORT} dumped_data.sql ${DB}

# Create the database in GREEN to receive the data from BLUE
./sql_execute.sh ${GREEN_HOST} ${GREEN_PORT} create_db.sql

# Import the data from the file into GREEN
./sql_execute.sh ${GREEN_HOST} ${GREEN_PORT} dumped_data.sql ${DB}

# Migrate the data! In this case, perform an "ALTER TABLE" changing the column "name" to "nickname"
./sql_execute.sh ${GREEN_HOST} ${GREEN_PORT} alter_table.sql ${DB}

# Optionally we could archive the dumped data at this point
