#!/bin/bash

set -ex
cd `dirname $0`

docker build -t malekvictor/space-box-writer:latest --target=app -f ../Dockerfile ..
#docker push malekvictor/space-box-writer:latest
