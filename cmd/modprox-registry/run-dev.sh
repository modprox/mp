#!/bin/bash

set -euo pipefail

go clean
go generate
go build

if [[ ${#} -eq 1 ]]; then
	configfile="${1}"
else
	configfile="../../hack/configs/registry-local.mysql.json"
fi

./modprox-registry ${configfile}
