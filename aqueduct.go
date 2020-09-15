package main

import (
	"log"
	"net/http"
	"time"

	"github.com/bdkiran/aqueduct/api"
)

func main() {

	router := api.NewRouter()
	port := ":8080"

	log.Printf("Server is launched at port %s", port)

	srv := &http.Server{
		Handler:      router,
		Addr:         port,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
