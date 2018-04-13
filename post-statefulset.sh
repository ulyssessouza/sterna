#!/usr/bin/env bash

set -e
set -xv

export BLUE_HOST="192.168.99.100"
export BLUE_PORT=30251

export GREEN_HOST=${BLUE_HOST}
export GREEN_PORT=30252

export DB="test"

pushd _examples/cockroachdb_examples_sql
    ./wait_clusters.sh
    ./migrate_data.sh
popd
