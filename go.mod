module github.com/modprox/modprox-proxy

require (
	github.com/boltdb/bolt v1.3.1
	github.com/gorilla/context v1.1.1 // indirect
	github.com/gorilla/mux v1.6.2
	github.com/modprox/libmodprox v0.0.0
	github.com/pkg/errors v0.8.0
	github.com/shoenig/atomicfs v0.1.1
	github.com/shoenig/toolkit v1.0.0
	github.com/stretchr/testify v1.2.2
)

replace github.com/modprox/libmodprox => ../libmodprox
