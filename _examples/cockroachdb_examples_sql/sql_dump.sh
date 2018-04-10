#!/bin/bash
cockroach dump test --insecure --host 192.168.99.100 -p 32186 > dump.sql