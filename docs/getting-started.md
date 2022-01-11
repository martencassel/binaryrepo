# Build

```bash
make clean
make build
```

# Run the server

The server is listening on address http://localalhost:8081

```
./binaryrepo
```

# Remote docker repos

There is a preconfigured remote repo with the name "docker-remote".
It's available on this URL:

http://localhost:8081/repo/docker-remote/v2

Docker is not able to access this path dute to limitations in the docker client.
Docker clients needs a /v2/ int the base url.

In order to solve this, a reverse proxy (nginx) can be used to make this repo available.

http://docker-remote.example.com -> http://localhost:8081/repo/docker-remote/v2

# Start the nginx

In order to start nginx, a local-ca is needed (see the setup-certs.sh script).

cd tools/docker
./start-nginx.sh

Then update /etc/hosts

```
127.0.0.1   localhost localhost.localdomain localhost4 localhost4.localdomain4
::1         localhost localhost.localdomain localhost6 localhost6.localdomain6
127.0.0.1   docker-remote.example.com
```

# Pull images from docker hub

unset http_proxy
unset https_proxy

sudo rm -rf /tmp/filestore/*

docker rmi -f redis
docker rmi -f docker-remote.example.com/redis
time docker image pull docker-remote.example.com/redis

docker rmi -f project.example.com/postgres
time docker image pull docker-remote.example.com/postgres

docker rmi -f postgres
docker rmi -f project.example.com/postgres
time docker image pull docker-remote.example.com/postgres

docker rmi -f postgres
docker rmi -f project.example.com/postgres
time docker image pull postgres

# Debugging network traffic

unset http_proxy
unset https_proxy
./mitmproxy -k

# Test connectivity

unset http_proxy
unset https_proxy
curl -vv http://localhost:8081/repo/docker-remote/v2;