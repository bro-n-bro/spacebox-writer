#!/bin/bash

set -ex
cd `dirname $0`

docker buildx build --platform linux/amd64 -t malekvictor/space-box-writer:latest --target=app ../