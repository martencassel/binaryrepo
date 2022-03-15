# About this project

Binaryrepo is a binary repository manager. It can manage binary packages of the following types:

* docker

It supports the following repository types (for docker)

* Remotes
* Local

## Background (Remotes)
Performance of package managers (Docker, Helm, Go etc) can be significantly improved by reusing previously fetched resources from the internet to a shared cache server.

Binaryrepo is a shared cache server that store responses to be reused by more than one user.
Multiple users may need to download a certain package from a repository on the internet.
By setting up a shared cache on the local network it may serve many users so that popular packages are
resused a number of times, reducing network traffic and latency.
## Features

* Remote docker repos, ie. caching proxy for remote repos in docker registries
* Local docker repos (A Docker Registry v2)
* Store files using checksums.
* Manage repos using an API (TODO)

## Demo

In the demo below, you can see that the the download time of postgres:latest will be reduced

[![asciicast](https://asciinema.org/a/1bHV8eIiAFO4t2G5Azx8HrqLs.svg)](https://asciinema.org/a/1bHV8eIiAFO4t2G5Azx8HrqLs)


# Getting started

Add a host entry for the reverse proxy:
```bash
127.0.0.1 api.binaryrepo.local docker-local.binaryrepo.local docker-remote.binaryrepo.local
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

# Usage

## Local docker repository

Create a local docker repository
```bash
curl --header "Content-Type: application/json" \
  --request POST \
  --data '{"name":"docker-local", "repo_type":"local","package_type":"docker"'\
  https://api.binaryrepo.local/api/repository
```

Push an image to the local docker repo:
```bash
docker image pull redis:latest
docker image tag redis:latest docker-local.binaryrepo.local/redis:latest
docker image push docker-local.binaryrepo.local/redis:latest
```

## Remote docker repository

Create a remote docker repository for docker hub,

```bash
curl --header "Content-Type: application/json" \
  --request POST \
  --data '{"name":"docker-remote", "repo_type":"remote","package_type":"remote","remote_url":"https://registry-1.docker.io"}' \
  https://api.binaryrepo.local/api/repository
```

If you have a docker hub account, use define using credentials,

```bash
curl --header "Content-Type: application/json" \
  --request POST \
  --data '{"name":"docker-remote","repo_type":"remote","package_type":"remote","username": "hub-user", "password": "hub-password", remote_url":"https://registry-1.docker.io"}' \
  https://api.binaryrepo.local/api/repository
```

Pull an image from the remote docker repo:

```bash
docker image rmi docker-remote.binaryrepo.local/redis:latest redis:latest
time docker image pull docker-remote.binaryrepo.local/redis:latest
```

Clear the local docker cache,

```bash
docker image rmi docker-remote.binaryrepo.local/redis:latest redis:latest
```

Then pull again, now from using the binaryrepo server cache,

```bash
time docker image pull docker-remote.binaryrepo.local/redis:latest
```