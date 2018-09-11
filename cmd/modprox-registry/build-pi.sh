#!/bin/bash

set -euo pipefail

go clean
go generate
GOOS=linux GOARCH=arm GOARM=7 go build

