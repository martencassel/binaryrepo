#!/bin/bash

docker-compose build

docker-compose down

docker-compose up -d

docker exec -it binaryrepo_binaryrepo_1 make binaryrepo

docker exec -it binaryrepo_binaryrepo_1 ./build/binaryrepo run

docker exec -it binaryrepo_binaryrepo_1 ps aux

docker exec -it binaryrepo_binaryrepo_1 /bin/sh


dlv attach `pidof binaryrepo` --listen=:2345 --headless --api-version=2 --log

docker exec -it binaryrepo_binaryrepo_1 dlv debug -l 0.0.0.0:2345 --headless=true --log=true ./cmd/binary-repo -- run


