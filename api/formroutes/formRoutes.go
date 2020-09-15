package formroutes

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/bdkiran/aqueduct/api/response"
	"github.com/gorilla/mux"
)

//FormHandler handles post request that contain form data. Currently processes JSON and url-endcoded payloads
func FormHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Form handler called")

	//Get the api key from the request
	vars := mux.Vars(r)
	log.Printf("key: %s", vars["key"])
	requestAPIKey := vars["key"]
	//Verify that the key is valid
	if !verifyKey(requestAPIKey) {
		response.SendFailResponse(w, "Invalid API key.")
		return
	}

	//Get the content-type header of the request
	content := r.Header.Get("Content-Type")
	var formContent []byte
	var err error

	//add form-data processing
	if content == "application/json" {
		formContent, err = parseJSONEndcodedData(r)
		if err != nil {
			response.SendFailResponse(w, "Unable to parse json data.")
			return
		}
	} else if content == "application/x-www-form-urlencoded" {
		formContent, err = parseURLEncodedData(r)
		if err != nil {
			response.SendFailResponse(w, "Unable to parse url encoded data.")
			return
		}
	} else {
		log.Printf("Unable to parse format %s", content)
		response.SendFailResponse(w, "Unable to parse format")
		return
	}
	log.Println(string(formContent))
	response.SendSuccessResponse(w, "Great success")
}

func parseURLEncodedData(r *http.Request) ([]byte, error) {
	err := r.ParseForm()
	if err != nil {
		log.Println("Unable to parse the form payload.")
		return nil, err
	}
	//Make a map of url encoded data
	mapData := make(map[string]string)
	for key, value := range r.Form {
		log.Printf("key: %s, value: %s", key, value[0])
		mapData[key] = value[0]
	}
	//Convert the map into a json object
	jsnData, err := json.Marshal(mapData)
	if err != nil {
		return nil, err
	}
	return jsnData, nil
}

func parseJSONEndcodedData(r *http.Request) ([]byte, error) {
	jsn, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return jsn, nil
}
