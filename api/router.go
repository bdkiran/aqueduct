package api

import (
	"log"
	"net/http"

	chatapi "github.com/bdkiran/aqueduct/api/chatroutes"
	"github.com/bdkiran/aqueduct/api/formroutes"
	"github.com/bdkiran/aqueduct/api/users"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

//NewRouter registers all routers and returns a router.
func NewRouter() http.Handler {
	router := mux.NewRouter()
	router.HandleFunc("/", homeHandler).Methods("GET")
	router.HandleFunc("/key", formroutes.RequestKey).Methods("GET")
	router.HandleFunc("/f/{key}", formroutes.FormHandler).Methods("POST")

	//Serve static thank you page...
	dir := ("./pages/")
	router.PathPrefix("/thanks/").Handler(http.StripPrefix("/thanks/", http.FileServer(http.Dir(dir))))

	//Register the user routes
	users.Register(router)
	chatapi.RegisterConnection(router)

	//Only certain routes should be CORS like the form and chat routes?
	headersOk := handlers.AllowedHeaders([]string{"Accept", "Accept-Language", "X-Requested-With", "Content-Type", "Authorization"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})
	returnCors := handlers.CORS(headersOk, originsOk, methodsOk)(router)
	return returnCors
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Home handler called")
}
