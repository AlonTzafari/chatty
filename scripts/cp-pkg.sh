#!/bin/bash

mkdir -p .pkg
find ./apps -name 'package.json' -type f -exec cp --parents "{}" ./.pkg \;
find ./apps -name 'go.mod' -type f -exec cp --parents "{}" ./.pkg \;
find ./libs -name 'package.json' -type f -exec cp --parents "{}" ./.pkg \;
find ./apps -name 'go.sum' -type f -exec cp --parents "{}" ./.pkg \;
cp package.json .pkg/package.json
cp package-lock.json .pkg/package-lock.json