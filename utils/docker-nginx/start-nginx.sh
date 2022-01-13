#!/bin/bash

unset http_proxy 
unset https_proxy

# Remove any old instances
docker rm -f $(docker ps -aq)
sudo rm -rf /tmp/nginx.conf/

# Copy the configuration file to /tmp/nginx.conf
sudo cp ./nginx.conf /tmp/nginx.conf

# Start nginx to proxy againt the remote repo
docker image pull nginx:latest
docker run --net=host --name docker-remote-proxy \
    -v ~/local-ca/certs:/certs \
    -v /tmp/nginx.conf:/etc/nginx/nginx.conf:ro \
    -d nginx

# Connect directly to the binarrepo server
unset http_proxy https_proxy
curl -k -vv http://localhost:8081/repo/docker-remote/v2

# Connect via nginx
unset http_proxy https_proxy
curl -k -v https://docker-remote.example.com/v2/

