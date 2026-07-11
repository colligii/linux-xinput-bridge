#!/bin/bash

set -e

go build -o xbox-driver xbox.go
go build -o controller-reader read_controller.go

echo "Go binaries built successfully"