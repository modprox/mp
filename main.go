package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/modprox/modprox-proxy/web"
)

func main() {
	fmt.Println("starting the modprox-proxy service")

	router := web.NewRouter()
	if err := http.ListenAndServe(":10001", router); err != nil {
		log.Fatalf("failed to listen and serve forever %v", err)
	}
}
