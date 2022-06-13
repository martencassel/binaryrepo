#!/bin/bash

# Start the server

# Create two remote repos using the API

curl --header "Content-Type: application/json" \
  --request POST \
  --data '{"name":"alpine-main", "repo_type":"remote","package_type":"alpine","remote_url":"https://dl-cdn.alpinelinux.org/alpine/v3.14/main"}' \
  http://localhost:8080/api/repository

curl --header "Content-Type: application/json" \
  --request POST \
  --data '{"name":"alpine-community", "repo_type":"remote","package_type":"alpine","remote_url":"https://dl-cdn.alpinelinux.org/alpine/v3.14/community"}' \
  http://localhost:8080/api/repository

# List the registered repos

curl -X GET http://localhost:8080/api/repository
{
    "repositories": [
        {
          "name": "alpine-main",
          "repo_type": "remote",
          "package_type": "alpine",
          "remote_url": "https://dl-cdn.alpinelinux.org/alpine/v3.14/main"
        },
        {
          "name": "alpine-community",
          "repo_type": "remote",
          "package_type": "alpine",
          "remote_url": "https://dl-cdn.alpinelinux.org/alpine/v3.14/community"
        }
    ]    
}


# Provision alpine packages from the remote repos above

cat << EOF > /tmp/repos
http://localhost:8080/apk-remote-main
http://localhost:8080/apk-remote-community
EOF

docker run --name=alpine -d -v /host:/tmp alpine sleep infinity
docker cp /tmp/repos alpine:/etc/apk/repositories
docker exec -it alpine cat /etc/apk/repositories
docker exec -it alpine apk update
docker exec -it alpine apk fetch vim

# Protocol

# apk update            fetch http://localhost:8080/alpine-remote/x86_64/APKINDEX.tar.gz
# apk fetch             fetch http://localhost:8080/alpine-remote/x86_64/vim-7.4.tar.gz
# apk add vim           fetch all dependencies



