#!/bin/bash

set -euo pipefail
set -x

go generate
go build

if [ ${#} -eq 0 ]; then
	./modprox-proxy ../../hack/configs/proxy-local.json
else
	./modprox-proxy "${1}"
fi
