#!/bin/bash

set -ex
cd `dirname $0`

docker buildx build --platform linux/amd64 -t malekvictor/space-box-writer:latest -f ../Dockerfile-amd --target=app ../
#docker buildx create --use desktop-linux
#docker buildx build --platform linux/arm64 -t malekvictor/space-box-writer:0.0.4 --load --target=app ../

