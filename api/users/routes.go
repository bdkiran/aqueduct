package users

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/bdkiran/aqueduct/api/response"

	"github.com/bdkiran/aqueduct/persist"
	"github.com/gorilla/mux"
)

//Register handles registering of all the user routers
func Register(r *mux.Router) {
	r.HandleFunc("/user/{id}", getUser).Methods("GET")
	r.HandleFunc("/user", registerUser).Methods("POST")
}

func getUser(w http.ResponseWriter, r *http.Request) {
	variables := mux.Vars(r)
	userID := variables["id"]
	log.Println(userID)
	user, err := persist.GetUserByID(userID)
	if err != nil {
		//send 400 response? does not match object
		response.SendFailResponse(w, user)
		return
	}
	response.SendSuccessResponse(w, user)
}

func registerUser(w http.ResponseWriter, r *http.Request) {
	var u persist.User
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&u)
	if err != nil {
		//send 400 response? does not match object
		response.SendFailResponse(w, "Error with the request object")
		return
	}
	responseUser, err := persist.InsertUser(u)
	if err != nil {
		//send 500 response? issue with insertion or database
		response.SendErrorResponse(w, nil, "Error inserting user")
		return
	}
	//Send 200 successful..
	response.SendSuccessResponse(w, responseUser)
}
