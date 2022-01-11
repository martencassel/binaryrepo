#!/bin/bash

unset http_proxy
unset https_proxy

docker rm -f $(docker ps -aq)
sudo rm -rf /tmp/nginx.conf/

sudo cp ./nginx.conf /tmp/nginx.conf

docker image pull nginx:latest

docker run --net=host --name my-custom-nginx-container \
    -v ~/local-ca/certs:/certs \
    -v /tmp/nginx.conf:/etc/nginx/nginx.conf:ro \
    -d nginx

