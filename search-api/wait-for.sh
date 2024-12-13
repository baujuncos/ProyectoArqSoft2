#!/bin/bash
set -e

HOST="$1"
PORT="$2"

while ! nc -z "$HOST" "$PORT"; do
  echo "Esperando a que $HOST:$PORT esté listo..."
  sleep 1
done

echo "$HOST:$PORT está listo!"
exec "$@"

