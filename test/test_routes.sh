#!/bin/bash

curl -I http://localhost:8081/repo/docker-remote/v2/alpine/manifests/latest
curl -L -X GET http://localhost:8081/repo/docker-remote/v2/library/alpine/blobs/sha256:c059bfaa849c4d8e4aecaeb3a10c2d9b3d85f5165c66ad3a4d937758128c4d18  

curl -I http://localhost:8081/repo/docker-remote/v2/alpine/manifests/sha256:f446e3fb1a01cc3e6199cbce8525d8b0ee3f2c0c4dd74ac0820bed19bb706335
curl -I http://localhost:8081/repo/docker-remote/v2/library/alpine/manifests/sha256:f446e3fb1a01cc3e6199cbce8525d8b0ee3f2c0c4dd74ac0820bed19bb706335
curl -I http://localhost:8081/repo/docker-not-found/v2/library/alpine/manifests/sha256:f446e3fb1a01cc3e6199cbce8525d8b0ee3f2c0c4dd74ac0820bed19bb706335

curl -I -X GET http://localhost:8081/repo/docker-remote/v2/alpine/manifests/sha256:f446e3fb1a01cc3e6199cbce8525d8b0ee3f2c0c4dd74ac0820bed19bb706335
curl -I -X GET http://localhost:8081/repo/docker-remote/v2/library/alpine/manifests/sha256:f446e3fb1a01cc3e6199cbce8525d8b0ee3f2c0c4dd74ac0820bed19bb706335
curl -I -X GET http://localhost:8081/repo/docker-not-found/v2/library/alpine/manifests/sha256:f446e3fb1a01cc3e6199cbce8525d8b0ee3f2c0c4dd74ac0820bed19bb706335

curl -I -X GET http://localhost:8081/repo/docker-remote/v2/alpine/blobs/sha256:f446e3fb1a01cc3e6199cbce8525d8b0ee3f2c0c4dd74ac0820bed19bb706335
curl -I -X GET http://localhost:8081/repo/docker-remote/v2/library/alpine/blobs/sha256:f446e3fb1a01cc3e6199cbce8525d8b0ee3f2c0c4dd74ac0820bed19bb706335
curl -I -X GET http://localhost:8081/repo/docker-not-found/v2/library/alpine/blobs/sha256:f446e3fb1a01cc3e6199cbce8525d8b0ee3f2c0c4dd74ac0820bed19bb706335
