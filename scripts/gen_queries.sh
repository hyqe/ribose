#!/bin/bash
set -ex
docker run \
    --rm \
    -v $(pwd):/src \
    -w /src \
    kjconroy/sqlc \
        generate