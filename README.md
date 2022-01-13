# About this project
## Motivation

Learn how various package managers work, and how content are downloaded from the internet.
Improve the way packages are downloaded and managed by implementing a all-in-one binary repository manager
that can serve multiple packages types, local access or remote access (proxy cache).

Downloading docker images from docker hub can be slow.

Solution:

Binaryrepo serves a caching proxy for any remote registry that implements Docker Registry v2 API.
Requests are served from its local cache (the file storage), if not available in the cache, requests
are forwarded to the remote registry, and resources are fetched and saved to the cache for future requests.

# Future plans

I plan to implement a local Docker Registry v2, that can be used together with caching proxy functionality.
I may also implement support for other package types such as Go modules and Helm packages.

# binary-repo

Please add a small section about the problem and how this project solves


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
