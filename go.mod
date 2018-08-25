module github.com/modprox/modprox-registry

require (
	github.com/go-sql-driver/mysql v1.4.0
	github.com/gorilla/context v1.1.1 // indirect
	github.com/gorilla/mux v1.6.2
	github.com/modprox/libmodprox v0.0.0
	github.com/pkg/errors v0.8.0
	google.golang.org/appengine v1.1.0 // indirect
)

replace github.com/modprox/libmodprox => ../libmodprox
