package main

import (
	"log"
	"net/http"

	"github.com/go-sql-driver/mysql"
	"github.com/modprox/modprox-registry/internal/repositories"
	"github.com/modprox/modprox-registry/internal/web"
)

func main() {
	dsn := mysqlDSN()
	db, err := repositories.Connect(dsn)
	if err != nil {
		log.Fatalf("failed to connect to mysql database: %v", err)
	}

	store, err := repositories.New(db)
	if err != nil {
		log.Fatalf("failed to create registry database: %v", err)
	}

	router := web.NewRouter(store)
	if err := http.ListenAndServe(":8000", router); err != nil {
		log.Fatalf("failed to listen and serve forever %v", err)
	}
}

func mysqlDSN() mysql.Config {
	return mysql.Config{
		User:                 "docker",
		Passwd:               "docker",
		Addr:                 "localhost:3306",
		DBName:               "modproxdb",
		AllowNativePasswords: true,
	}
}
