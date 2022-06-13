# About this project

Binaryrepo is a binary repository manager. It can manage binary packages of the following types:

* docker

It supports the following repository types (for docker)

* Remotes
* Local

Accessing docker repos from binaryrepo server requires a reverse proxy in front
of the binaryrepo server. See docker-compose.yaml

## Background (Remotes)
Performance of package managers (Docker, Helm, Go etc) can be significantly improved by reusing previously fetched resources from the internet to a shared cache server.

Binaryrepo is a shared cache server that store responses to be reused by more than one user.
Multiple users may need to download a certain package from a repository on the internet.
By setting up a shared cache on the local network it may serve many users so that popular packages are
resused a number of times, reducing network traffic and latency.
## Features

* Local docker registry (Docker Registry v2) with basic auth support
* Checksum based file management (filestore)
* Metadata management (file metdata management) in postgres.
* User store in postgres (password stored as hash)
* Postgres support (schema management, data access layer) 
* Installation using docker-compose
* Development tooling (VSCode devcontainer, docker-compose)

In-progress

* Docker Remote proxying
* Management REST API (demo UX) 

Todo

* Docker Registry V2 input validation (URL parameters etc)
* Support proxying of docker registries (Remote docker repos)
* Support helm packaging (remotes, locals)
* Support generic http file management (remotes, local
* Support user repository authorization layer (Access control for repositories)

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

Start binaryrepo with nginx and postgres
```bash
docker-compose up -d
```

# Usage

## Local docker repository

Login with default admin user
```bash
docker login api.binaryrepo.local -u admin
Password: admin
```

Create a local docker repository
```bash
curl -vv -u admin:admin -k --header "Content-Type: application/json" \
  --request POST \
  --data '{"name":"docker-local", "repo_type":"local","pkg_type":"docker"}'\
  https://api.binaryrepo.local/api/repo
```

List repos
```bash
curl -u admin:admin -k https://api.binaryrepo.local/api/repo|jq
```

Push an image to the docker-local repo
```bash
docker image pull alpine:latest
docker image tag alpine:latest docker-local.binaryrepo.local/alpine:latest
docker image tag alpine:latest docker-local.binaryrepo.local/alpine:tag1
docker image tag alpine:latest docker-local.binaryrepo.local/alpine:tag2
docker image tag alpine:latest docker-local.binaryrepo.local/alpine:tag3
docker image tag alpine:latest docker-local.binaryrepo.local/alpine:tag4
docker image tag alpine:latest docker-local.binaryrepo.local/alpine:tag5
docker image push docker-local.binaryrepo.local/alpine -a
```

## Remote docker repository

Create a remote docker repository for docker hub,

```bash
curl -v -k --header "Content-Type: application/json" \
  --request POST \
  --data '{"name":"docker-remote", "repo_type":"remote","package_type":"remote","remote_url":"https://registry-1.docker.io"}' \
  https://api.binaryrepo.local/api/repo
```

If you have a docker hub account, use define using credentials,

```bash
curl -v -k --header "Content-Type: application/json" \
  --request POST \
  --data '{"name":"docker-remote","repo_type":"remote","package_type":"remote","username": "hub-user", "password": "hub-password", remote_url":"https://registry-1.docker.io"}' \
  https://api.binaryrepo.local/api/repo
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
this operation will be faster than the previous one (the cache was empty).

```bash
time docker image pull docker-remote.binaryrepo.local/redis:latest
```
