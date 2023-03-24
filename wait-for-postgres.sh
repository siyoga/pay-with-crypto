#!/bin/sh
# wait-for-postgres.sh

set -e

host="$1"
shift
cmd="$@"

until PGPASSWORD=$POSTGRES_PASSWORD psql -h "$host" -U "admin" -d "main" -c '\q'; do
  >&2 echo "Postgres in unavaliable - sleeping"
  sleep 1
done

>&2 echo "Postgres is up - executing command"
exec $cmd