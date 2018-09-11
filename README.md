# modprox - The Go Module Proxy
modprox is a Go Module Proxy focused on the internal hosting use case

documentation @ [modprox.org](https://modprox.org)

[![Go Report Card](https://goreportcard.com/badge/github.com/modprox/mp)](https://goreportcard.com/report/github.com/modprox/mp) 
[![Build Status](https://travis-ci.org/modprox/mp.svg?branch=master)](https://travis-ci.org/modprox/mp) 
[![GoDoc](https://godoc.org/github.com/modprox/mp?status.svg)](https://godoc.org/github.com/modprox/mp) 
[![License](https://img.shields.io/github/license/modprox/mp.svg?style=flat-square)](LICENSE)


#### Project Management

- Issue [tracker](https://github.com/modprox/mp/issues)
- for the registry, prefix issues with "registry:"
- for the proxy, prefix issues with "proxy:"
- for the library packages (pkg/), prefix issues with "lib:"

#### Setting up modprox in your environment
For setting up your own instances of the modprox components, check out the
extensive documentation on [modprox.org](https://modprox.org/#starting)

#### Hacking on the Registry

The registry needs a persistent store, and for local development we have a docker image
with PostgreSQL setup to automatically create tables and users. To make things super simple, in
the `hack/` directory there is a `docker-compose` file already configured to setup the basic
containers needed for local developemnt. Simply run
```bash
$ docker-compose up
```
in the `hack/` directory to get them going. Also in the `hack/` directory is a script for
connecting to the MySQL that is running in the docker container, for ease of poking around.
```bash
$ compose up
Recreating modprox-postgres ... 
Recreating modprox-postgres ... done
Attaching to modprox-postgres
modprox-postgres | 2018-09-11 02:19:14.322 UTC [1] LOG:  listening on IPv4 address "0.0.0.0", port 5432
modprox-postgres | 2018-09-11 02:19:14.322 UTC [1] LOG:  listening on IPv6 address "::", port 5432
modprox-postgres | 2018-09-11 02:19:14.339 UTC [1] LOG:  listening on Unix socket "/var/run/postgresql/.s.PGSQL.5432"
modprox-postgres | 2018-09-11 02:19:14.370 UTC [22] LOG:  database system was shut down at 2018-09-11 02:19:12 UTC
modprox-postgres | 2018-09-11 02:19:14.381 UTC [1] LOG:  database system is ready to accept connections
```

Also in the `hack/` directory are some sample configuration files. By default, the included `run-dev.sh`
script will use the `hack/configs/registry-local.postgres.json` file, which works well with the included
`docker-compose.yaml` file.

#### Hacking on the Proxy

The proxy component is more simple than the registry in that it does not connect to anything (other than the registry
itself). It does however maintain its data-store of downloaded modules on disk, and by default it saves modules in the `/tmp`
directory.
