#!/bin/bash

set -euo pipefail


PGPASSWORD=docker psql \
	--host 127.0.0.1 \
	--port 5432 \
	--username docker \
	modproxdb

