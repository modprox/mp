# modprox-registry
registry component of modprox

#### Project Management
Issue [tracker](https://github.com/modprox/modprox-registry/issues) (prefix with `registry:`)

#### Setting up the Registry
For setting up your own instance of the modprox-registry,
see the documentation on [modprox.org](https://modprox.org/#starting)

#### Hacking on the Registry

The registry needs a persistent store, and for local development we have a docker image
with MySQL setup to automatically create tables and users. To make things super simple, in
the `hack/` directory there is a `docker-compose` file already configured to setup the basic
containers needed for local developemnt. Simply run
```bash
$ (modprox-registry/hack) docker-compose up
```
in the `hack/` directory to get them going. Also in the `hack/` directory is a script for
connecting to the MySQL that is running in the docker container, for ease of poking around.
```bash
$ (modprox-registry/hack) ./connect-mysql.sh
Welcome to the MySQL monitor.  Commands end with ; or \g.
Your MySQL connection id is 2
Server version: 5.7.22 MySQL Community Server (GPL)
```

Also in the `hack/` directory is `configs/local.json`. This file is used by the `run-dev.sh`
script that sits in the root of the repository. The `local.json` file is just some default
configuration settings for the registry that are meant to work with the provide docker images
for local development.
