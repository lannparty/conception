#!/bin/bash

for i in \
  github.com/docker/docker/api/types \
  github.com/docker/docker/api/types/container \
  github.com/docker/docker/client \
  github.com/docker/docker/pkg/stdcopy \
  golang.org/x/net/context
do
  go get $i
done
go build
