#!/bin/bash

# Unset proxy variables
unset http_proxy
unset https_proxy

# Pull images through the remote proxy repo

# Clear the local docker cache.
docker rmi -f docker-remote.example.com/redis:latest
docker rmi -f docker-remote.example.com/alpine:latest
docker rmi -f docker-remote.example.com/postgres:latest

time docker pull docker-remote.example.com/redis:latest
time docker pull docker-remote.example.com/alpine:latest
time docker pull docker-remote.example.com/postgres:latest

# Now the binaryrepo filestore is populated with image layers.

# Clear the local docker cache once again.
docker rmi -f docker-remote.example.com/redis:latest
docker rmi -f docker-remote.example.com/alpine:latest
docker rmi -f docker-remote.example.com/postgres:latest

# Now we will pull layers from binaryrepo's filestore cache.
time docker pull docker-remote.example.com/redis:latest
time docker pull docker-remote.example.com/alpine:latest
time docker pull docker-remote.example.com/postgres:latest

# Here goes the cache
tree /tmp/filestore

