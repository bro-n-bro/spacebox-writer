#!/bin/bash
set -ex
cd `dirname $0`
docker buildx build --platform linux/amd64 -t hexydev/spacebox-writer:latest --load -f ../Dockerfile-amd --target=app ../
