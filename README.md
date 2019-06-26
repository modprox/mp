# modprox - The Go Module Proxy
modprox is a Go Module Proxy focused on the internal hosting use case

documentation @ [modprox.org](https://modprox.org)

[![Go Report Card](https://goreportcard.com/badge/github.com/modprox/mp)](https://goreportcard.com/report/github.com/modprox/mp) 
[![Build Status](https://travis-ci.org/modprox/mp.svg?branch=master)](https://travis-ci.org/modprox/mp) 
[![GoDoc](https://godoc.org/github.com/modprox/mp?status.svg)](https://godoc.org/github.com/modprox/mp) 
[![License](https://img.shields.io/github/license/modprox/mp.svg?style=flat-square)](LICENSE)


#### Project Management

- [Issues](https://github.com/modprox/mp/issues)
- for the registry, prefix issues with `registry:`
- for the proxy, prefix issues with `proxy:`
- for the library packages (pkg/), prefix issues with `lib:`

#### Setting up modprox in your environment
For setting up your own instances of the modprox components, check out the
extensive documentation on [modprox.org](https://modprox.org/#starting)

#### Hacking on the Registry

The registry needs a persistent store, and for local development we have a docker image
with MySQL setup to automatically create tables and users. To make things super simple, in
the `hack/` directory there is a `docker-compose` file already configured to setup the basic
containers needed for local developemnt. Simply run
```bash
$ docker-compose up
```
in the `hack/` directory to get them going. Also in the `hack/` directory is a script for
connecting to the MySQL that is running in the docker container, for ease of poking around.
```bash
$ compose up
Starting modprox-fakeadog       ... done
Starting modprox-mysql-proxy    ... done
Starting modprox-mysql-registry ... done
Attaching to modprox-mysql-proxy, modprox-mysql-registry, modprox-fakeadog
modprox-mysql-registry | [Entrypoint] MySQL Docker Image 5.7.26-1.1.11
modprox-mysql-registry | [Entrypoint] Initializing database
modprox-mysql-proxy | [Entrypoint] MySQL Docker Image 5.7.26-1.1.11
modprox-mysql-proxy | [Entrypoint] Initializing database
modprox-mysql-proxy | [Entrypoint] Database initialized
modprox-fakeadog  | time="2019-06-26T18:19:14Z" level=info msg="listening on 0.0.0.0:8125"
modprox-mysql-registry | [Entrypoint] Database initialized
```

Also in the `hack/` directory are some sample configuration files. By default, the included `run-dev.sh`
script will use the `hack/configs/registry-local.mysql.json` file, which works well with the included
`docker-compose.yaml` file.

#### Hacking on the Proxy

The Proxy needs to persist its data-store of downloaded modules. It can be configured to either persist them to disk
or to MySQL.

##### local disk config
```json
"module_storage": {
  "data_path": "<disk path to store data>",
  "index_path": "<disk path to store boltdb index>",
  "tmp_path": "<disk path to store temporary files>"
}
```
Note that `data_path` and `tmp_path` should point to paths on the same filesystem.

##### MySQL config
```json
"module_db_storage": {
  "mysql": {
    "user": "docker",
    "password": "docker",
    "address": "localhost:3306",
    "database": "modproxdb-prox",
    "allow_native_passwords": true
  }
}
```
