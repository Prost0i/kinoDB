#!/bin/sh

export GOOSE_DRIVER=postgres
export GOOSE_DBSTRING="host=localhost port=5432 user=postgres dbname=kinodb sslmode=disable"

goose -dir migrations "$@"
