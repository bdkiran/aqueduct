package api

import (
	"log"
	"net/http"

	"github.com/bdkiran/aqueduct/api/formroutes"
	"github.com/gorilla/mux"
)

//NewRouter registers all routers and returns a router.
func NewRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/", homeHandler).Methods("GET")
	router.HandleFunc("/key", formroutes.RequestKey).Methods("GET")
	router.HandleFunc("/f/{key}", formroutes.FormHandler)
	return router
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Home handler called")
}
