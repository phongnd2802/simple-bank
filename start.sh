#!/bin/sh

set -e

echo "run db migration"

export GOOSE_DRIVER=postgres
export GOOSE_DBSTRING="user=root password=secret host=postgres port=5432 dbname=simple-bank sslmode=disable"
export GOOSE_MIGRATION_DIR=migration


ls -al
./goose -dir=$GOOSE_MIGRATION_DIR up

echo "start server"
exec "$@"