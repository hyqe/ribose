#!/bin/bash
# $1=POSTGRES_URL
# $2=up|down
# $3=STEP
#
# ./scripts/migrate.sh <POSTGRES_URL> <up|down|force> <STEP>
#
# https://github.com/golang-migrate/migrate/tree/v4.16.2/cmd/migrate#usage
set -ex
docker run \
    --rm \
    -v $(pwd)/internal/database/schema:/migrations \
    --network host \
    migrate/migrate \
        -path=/migrations/ \
        -database $1 \
        $2 $3 