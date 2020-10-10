package main

import (
	"log"
	"net/http"
	"time"

	"github.com/bdkiran/aqueduct/api"
	"github.com/bdkiran/aqueduct/persist"
)

func main() {
	router := api.NewRouter()
	port := ":8080"

	// mongoUsername := "admin"
	// mongoPassword := "aqueductpassword"
	// persist.InitMongoClient(mongoUsername, mongoPassword)

	log.Printf("Server is launched at port %s", port)
	srv := &http.Server{
		Handler:      router,
		Addr:         port,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
	defer persist.DisconnectClient()
}
