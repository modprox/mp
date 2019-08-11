# modprox - The Private Go Module Proxy
modprox is a Go Module Proxy focused on the private internal hosting use case

documentation @ [modprox.org](https://modprox.org)

[![Go Report Card](https://goreportcard.com/badge/github.com/modprox/mp)](https://goreportcard.com/report/github.com/modprox/mp) 
[![Build Status](https://travis-ci.org/modprox/mp.svg?branch=master)](https://travis-ci.org/modprox/mp) 
[![GoDoc](https://godoc.org/github.com/modprox/mp?status.svg)](https://godoc.org/github.com/modprox/mp)
[![NetflixOSS Lifecycle](https://img.shields.io/github.com/modprox/mp.svg)](OSSMETADATA)
[![License](https://img.shields.io/github/license/modprox/mp.svg?style=flat-square)](LICENSE)

# Project Overview

Module `oss.indeed.com/go/modprox` provides a solution for hosting an internal
Go Module Proxy that is capable of communicating with private authenticated
git repositories.

# Contributing

We welcome contributions! Feel free to help make `modprox` better.

#### Process

- We track bugs / features in [Issues](https://github.com/modprox/mp/issues)
- Open an issue and describe the desired feature / bug fix before making
changes. It's useful to get a second pair of eyes before investing development
effort.
- Make the change. If adding a new feature, remember to provide tests that
demonstrate the new feature works, including any error paths. If contributing
a bug fix, add tests that demonstrate the erroneous behavior is fixed.
- Open a pull request. Automated CI tests will run. If the tests fail, please
make changes to fix the behavior, and repeat until the tests pass.
- Once everything looks good, one of the indeedeng members will review the
PR and provide feedback.

### Setting up modprox in your environment
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

# Asking Questions

For technical questions about `modprox`, just file an issue in the GitHub tracker.

For questions about Open Source in Indeed Engineering, send us an email at
opensource@indeed.com

# Maintainers

The `oss.indeed.com/go/modprox` module is maintained by Indeed Engineering.

While we are always busy helping people get jobs, we will try to respond to
GitHub issues, pull requests, and questions within a couple of business days.

# Code of Conduct

`oss.indeed.com/go/modprox` is governed by the[Contributer Covenant v1.4.1](CODE_OF_CONDUCT.md)

For more information please contact opensource@indeed.com.

# License

The `oss.indeed.com/go/modprox` module is open source under the [BSD-3-Clause](LICENSE)
license.
