#!/bin/bash	

time golangci-lint run && time go test ./... -v --race
