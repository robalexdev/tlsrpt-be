#!/usr/bin/bash

set -e

docker compose --env-file prod.env --parallel=1 build
docker compose --env-file prod.env up -d --wait --remove-orphans

