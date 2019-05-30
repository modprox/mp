module github.com/modprox/mp

require (
	github.com/boltdb/bolt v1.3.1
	github.com/cactus/go-statsd-client v3.1.1+incompatible
	github.com/go-sql-driver/mysql v1.4.0
	github.com/gorilla/csrf v1.5.1
	github.com/gorilla/mux v1.6.2
	github.com/jinzhu/copier v0.0.0-20180308034124-7e38e58719c3
	github.com/lib/pq v1.0.0
	github.com/modprox/taggit v0.0.5
	github.com/pkg/errors v0.8.0
	github.com/shoenig/atomicfs v0.1.1
	github.com/shoenig/httplus v0.0.0
	github.com/shoenig/petrify/v4 v4.1.0
	github.com/shoenig/toolkit v1.0.0
	github.com/stretchr/testify v1.2.2
	google.golang.org/appengine v1.3.0 // indirect
)

exclude (
	// Version of petrify/v4 before v4.1.0 are bad, but libraries keep pulling them into go.sum.
	// This was caused by the go1.11.4 change where symlinks started being hashed differently.
	github.com/shoenig/petrify/v4 v4.0.2
	github.com/shoenig/petrify/v4 v4.0.3
	github.com/shoenig/petrify/v4 v4.0.4
	github.com/shoenig/petrify/v4 v4.0.5
)
