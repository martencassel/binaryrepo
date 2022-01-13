#!/bin/bash

# Unset proxy variables
unset http_proxy
unset https_proxy

# Start the reverse proxy
../utils/docker-nginx/start-nginx.sh

# Pull images through the remote proxy repo

docker pull docker-remote.example.com/redis:latest

docker pull docker-remote.example.com/alpine:latest

docker pull docker-remote.example.com/postgres:latest
