#!/bin/bash

unset http_proxy
unset https_proxy

sudo rm -rf /tmp/filestore/*

docker rmi -f redis
docker rmi -f docker-remote.example.com/redis

docker rmi -f project.example.com/postgres
time docker image pull docker-remote.example.com/postgres

docker rmi -f postgres
docker rmi -f docker-remote.example.com/postgres
time docker image pull docker-remote.example.com/postgres

docker rmi -f postgres
docker rmi -f docker-remote.example.com/postgres
time docker image pull postgres
