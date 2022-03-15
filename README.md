# About this project

## Problem

Performance of package managers (Docker, Helm, Go etc) can be significantly improved by reusing previously fetched resources from the internet to a shared cache server.

Binaryrepo is a shared cache server that store responses to be reused by more than one user.
Multiple users may need to download a certain package from a repository on the internet.
By setting up a shared cache on the local network it may serve many users so that popular packages are
resused a number of times, reducing network traffic and latency.
## Features

* Proxy caching images from docker hub using a docker hub login
* The local cache checksum based approach which reduces storage space by storing files only once.

## Demo

In the demo below, you can see that the the download time of postgres:latest will be reduced

[![asciicast](https://asciinema.org/a/1bHV8eIiAFO4t2G5Azx8HrqLs.svg)](https://asciinema.org/a/1bHV8eIiAFO4t2G5Azx8HrqLs)


# Getting started

Add a host entry for the reverse proxy:
```bash
127.0.0.1 binaryrepo.example.com docker-local.example.com docker-remote.example.com
```

Create certificates:
```bash
make setup-certs
```

Build:
```bash
docker-compose build
```

Start nginx and binaryrepo
```bash
docker-compose up -d
```

Create a local docker repository
```bash
curl --header "Content-Type: application/json" \
  --request POST \
  --data '{"name":"docker-local", "repo_type":"local","package_type":"docker"'\
  https://binaryrepo.example.com/api/repository
```

Create a remote docker repository
```bash
curl --header "Content-Type: application/json" \
  --request POST \
  --data '{"name":"docker-remote", "repo_type":"remote","package_type":"remote","remote_url":"https://registry-1.docker.io"}' \
  https://binaryrepo.example.com/api/repository
```

Push an image to the local docker repo:
```bash
docker image pull redis:latest
docker image tag redis:latest docker-local.example.com/redis:latest
docker image push docker-local.example.com/redis:latest
```

Pull an image from the remote docker repo:
```bash
docker image rmi docker-remote.example.com/redis:latest redis:latest
time docker image pull docker-remote.example.com/redis:latest
docker image rmi docker-remote.example.com/redis:latest redis:latest
time docker image pull docker-remote.example.com/redis:latest
```

## Code

The docker proxy

https://github.com/martencassel/binaryrepo/tree/main/pkg/docker/proxy

The docker registry client

https://github.com/martencassel/binaryrepo/tree/main/pkg/docker/client

The filestore
https://github.com/martencassel/binaryrepo/tree/main/pkg/filestore/fs
