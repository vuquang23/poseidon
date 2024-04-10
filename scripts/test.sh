#! /bin/bash

echo "Running test..."

go test -p 1 -v -coverpkg=./... -coverprofile=profile.cov ./... > test.log

go tool cover -func profile.cov
