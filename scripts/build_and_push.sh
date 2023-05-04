#!/bin/bash
set -ex
cd `dirname $0`
docker build -t hexydev/spacebox-writer:latest --target=app -f ../Dockerfile ..