#!/bin/bash
set -ex
pwd
docker run --rm -v $(pwd):/src -w /src kjconroy/sqlc generate