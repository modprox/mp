module github.com/modprox/modprox-registry

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/go-sql-driver/mysql v1.4.0
	github.com/gorilla/csrf v1.5.1
	github.com/gorilla/mux v1.6.2
	github.com/modprox/libmodprox v0.0.0
	github.com/pkg/errors v0.8.0
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/shoenig/petrify/v4 v4.0.2
	github.com/stretchr/testify v1.2.2
	google.golang.org/appengine v1.1.0 // indirect
)

replace github.com/modprox/libmodprox => ../libmodprox
