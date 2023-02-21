#!/bin/bash
set -ex
cd `dirname $0`
docker buildx build --platform linux/amd64 -t hexydev/spacebox-writer:0.0.9 --load -f ../Dockerfile-amd --target=app ../
