#!/bin/bash
docker build --tag votes-api -f ./votes-api/dockerfile ./votes-api
docker build --tag voter-api -f ./voter-api/dockerfile ./voter-api
docker build --tag poll-api -f ./poll-api/dockerfile ./poll-api

docker compose up