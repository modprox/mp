#!/bin/bash

set -euo pipefail

mysql \
    --protocol=tcp \
    --host=localhost \
    --port=3307 \
    --user=docker \
    --password=docker \
    --database=modproxdb-prox
