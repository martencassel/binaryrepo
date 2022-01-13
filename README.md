# About this project

This project tries to implement a binary artifact manager server that can serve software packages such as
Docker images, Helm charts or Go modules etc, either from a local storage or from remote sources (caching proxy).
All binary packages are stored in a "single-instance-store" using a checksum scheme.
## Motivation

Learn about various package managers and related protocols, and implement them using Golang.
## Current state

Consuming docker images from docker hub can be slow.
With a central server (binaryrepo), it can serve docker images from the internet
and cache them locally.

Binaryrepo serves a caching proxy for any remote registry that implements Docker Registry v2 API.
Requests are served from its local cache (the file storage), if not available in the cache, requests
are forwarded to the remote registry, and resources are fetched and saved to the cache for future requests.
##

Binaryrepo is installed locally using the local /etc/hosts name of docker-remote.example.com,
The below pull command will fetch redis images from docker hub, and in the process cache this image.

docker pull docker-remote.example.com/redis:latest

Any futher pulls will be served from the cache that binaryrepo hosts.

# Future plans

I plan to implement a local Docker Registry v2, that can be used together with caching proxy functionality.
I may also implement support for other package types such as Go modules and Helm packages.

# binary-repo

### Getting started

The following example will setup binaryrepo to be used
as a remote proxy cache of docker hub.

Docker pull command will access the remote repo through a local nginx container.

#### Prerequisites
1. Create certs
```bash
make setup-certs
ls ~/certs
```
2. Modify /etc/hosts
```bash
127.0.0.1 docker-remote.example.com
```
#### Start everything and run docker pull tests
```bash
make check-remote-pull
```
### Building

```bash
make build
```
### Working demo
Proxy docker images from docker hub

- Create certs under $HOME/certs/ for nginx, add add the certs to your local cert store.
- Add docker-remote.example.com to your /etc/hosts file
- Start binaryrepo server at http://localhost:8081/
- Start nginx, it will proxies requests from docker-remote.example.com to http://localhost:8081/repo/docker-remote
- Run docker image pull on some sample images, remove the local images, and then pull again, now from the local cache in binaryrepo
  (tree /tmp/filestore)

```bash
make setup-certs
make check-remote-pull
```
