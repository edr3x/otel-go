#!/usr/bin/env bash

if [[ $1 ]]; then
  docker compose down
  exit 0
fi

docker volume create grafana_data

mkdir -p tempo_data

chmod 777 tempo_data

docker compose up -d
