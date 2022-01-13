#!/bin/bash

# Unset proxy variables
unset http_proxy
unset https_proxy

# Pull images through the remote proxy repo

docker rmi -f docker-remote.example.com/redis:latest
docker rmi -f docker-remote.example.com/alpine:latest
docker rmi -f docker-remote.example.com/postgres:latest

docker pull docker-remote.example.com/redis:latest
docker pull docker-remote.example.com/alpine:latest
docker pull docker-remote.example.com/postgres:latest

docker rmi -f docker-remote.example.com/redis:latest
docker rmi -f docker-remote.example.com/alpine:latest
docker rmi -f docker-remote.example.com/postgres:latest

docker pull docker-remote.example.com/redis:latest
docker pull docker-remote.example.com/alpine:latest
docker pull docker-remote.example.com/postgres:latest

