## modprox-proxy
proxy component of modprox

#### Setting up the Proxy
For setting up your own instance of the modprox-proxy,
see the documentation on [modprox.org](https://modprox.org/#starting).

#### Hacking on the Proxy

The proxy makes use of the filesystem as a datastore. For local development,
it's typical to point the module storage at a directory in `/tmp`. The `hack/`
directory includes a sample configuration JSON file for local development called
`local.json`. In the project root directory, there is a file called `run-dev.sh`
which will generate, build then run modprox-proxy with the provided sample config.
The sample config assumes the registry was already launched likewise.
