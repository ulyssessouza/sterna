#!/bin/bash

HOST=${1}
PORT=${2}
FILE=${3}
DB=${4}

cockroach dump ${DB} --insecure --host ${HOST} -p ${PORT} > ${3}