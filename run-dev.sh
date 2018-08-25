#!/bin/bash

set -euo pipefail

go clean
go generate
go build
./modprox-proxy hack/configs/local.json

