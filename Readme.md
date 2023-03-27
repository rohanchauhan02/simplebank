## Pull postgress image from docker

docker pull postgress:12-alpine

## Port mapping

docker run --name postgress12 - p 5432:5432 -e POSTGRES_USER=root -ePOSTGRES_PASSWORD=root -d postgres:12-alpine

## Run command in cotainer

docker exec -it postgres12 psql -U root

## view cotainer logs

docker logs postgres12

## create db using docker and shell

docker exec -it postgres12 /bin/sh
createdb --username=root --owner=root simple_bank
psql simple_bank
dropdb simple_bank

docker exec -it postgres12 createdb --username=root --owner=root simple_bank

## Install golang-migrate

brew install golang-migrate

## Create migration file

migrate create -ext sql  -dir db/migration -seq init_schema


## using makefile

To run postgress run command

1. make postgres

To create db run

2. make createdb

To run migration

3. make migrateup


## installation SQLC

using homebrew
1. brew install sqlc

using docker

1. docker pull kjconroy/sqlc

To run sqlc

2. docker run --rm -v $(pwd):/src -w /src kjconroy/sqlc generate



