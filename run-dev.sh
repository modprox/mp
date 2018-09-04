#!/bin/bash

set -euo pipefail

go clean
go generate
go build

if [ ${#} -eq 0 ]; then
	./modprox-proxy hack/configs/local.json
else
	./modprox-proxy "${1}"
fi
