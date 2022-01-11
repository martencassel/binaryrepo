#!/bin/bash

# Create a remote repo using the API
curl --header "Content-Type: application/json" \
  --request POST \
  --data '{"name":"helm-remote", "repo_type":"remote","package_type":"helm","remote_url":"https://charts.bitnami.com/bitnami"}' \
  http://localhost:8080/api/repository

# List the registered repos

curl -X GET http://localhost:8080/api/repository
{
    "repositories": [
        {
          "name": "helm-remote",
          "repo_type": "remote",
          "package_type": "helm",
          "remote_url": "https://charts.bitnami.com/bitnami"
        }
    ]    
}

# Helm repo add
helm repo add bitnami http://localhost:8080/helm-remote
helm repo update # Server fetches index.yaml from remote repo and rewrites it locally
helm fetch bitnami/nginx # Client fetches the rewritten index.yaml and fetches the package from the server