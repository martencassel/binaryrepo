#!/bin/bash

unset http_proxy
unset https_proxy

# Remove any old instances
docker rm -f $(docker ps -aq)

# Start nginx to proxy againt the remote repo
docker image pull nginx:latest

docker run --net=host --name docker-remote-proxy \
    -v ~/local-ca/certs:/certs \
    -v $(pwd)/nginx.conf:/etc/nginx/nginx.conf:ro \
    -d nginx
