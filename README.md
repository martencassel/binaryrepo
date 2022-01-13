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
