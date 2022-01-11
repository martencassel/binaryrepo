
# Start the server

# Create a local repo using the API

curl --header "Content-Type: application/json" \
  --request POST \
  --data '{"name":"alpine-local", "repo_type":"local","package_type":"alpine","remote_url":"https://dl-cdn.alpinelinux.org/alpine/v3.14/main"}' \
  http://localhost:8080/api/repository

# Upload some packages

# Calculate metadata
curl --header "Content-Type: application/json" \
  --request POST \
  http://localhost:8080/api/repository/alpine-local/reindex

# Servers perform apk index generation in the background

# Provision alpine packages from the remote repos above

cat << EOF >> /tmp/ repos
https://dl-cdn.alpinelinux.org/alpine/v3.14/main
https://dl-cdn.alpinelinux.org/alpine/v3.14/community
http://localhost:8080/apk-local/alpine/v3.14/local
EOF

docker run --name=alpine -d -v /host:/tmp alpine sleep infinity
docker cp /tmp/repos alpine:/etc/apk/repositories
docker exec -it alpine cat /etc/apk/repositories
docker exec -it alpine apk update
docker exec -it alpine apk fetch <package-name>
