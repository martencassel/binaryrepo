#!/bin/bash

docker-compose build

docker-compose down

docker-compose up -d

docker exec -it binaryrepo_binaryrepo_1 make binaryrepo


