#!/bin/sh

export GOOSE_DRIVER=postgres
export GOOSE_DBSTRING="user=postgres dbname=kinodb sslmode=disable"

goose -dir migrations "$@"
