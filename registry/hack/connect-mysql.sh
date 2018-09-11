#!/bin/bash

set -euo pipefail

mysql \
    --protocol=tcp \
    --host=localhost \
    --port=3306 \
    --user=docker \
    --password=docker \
    --database=modproxdb
