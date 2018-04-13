#!/usr/bin/env bash
cockroach sql --insecure --host 192.168.99.100 --port 30252 --database test -e "select * from test_tb;"
