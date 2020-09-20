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
//If this solution is to truely work, we need to handle: honeypot processing, custom-url redirect
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
	//Check honeypot field
	isSpam := honeyPotCheck(formContent)
	if isSpam {
		log.Println("Form is spam.....")
		//Dont send email to user
		//Dont store form
	}

	//If honeypot fails, lets not do this step??
	//check redirect field
	redirectURL := checkRedirectURL(formContent)
	//If no redirect is specified just use the generic thank you page
	if redirectURL == "" {
		redirectURL = "/thanks/"
	}

	//send email

	//Create thank you redirect page.
	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
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
		//log.Printf("key: %s, value: %s", key, value[0])
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

func honeyPotCheck(formJSONObj []byte) bool {
	var formMap map[string]interface{}
	err := json.Unmarshal(formJSONObj, &formMap)
	//this should never happen....
	if err != nil {
		log.Println("Big Error")
	}
	//Check if the honeypot field is contained in the form.
	if honeyValue, ok := formMap["honeypot"]; ok {
		if honeyValue != "" {
			return true
		}
		return false
	}
	return false
}

func checkRedirectURL(formJSONObj []byte) string {
	var formMap map[string]interface{}
	err := json.Unmarshal(formJSONObj, &formMap)
	//this should never happen....
	if err != nil {
		log.Println("Big Error")
	}
	//Check if the honeypot field is contained in the form.
	if redirectValue, ok := formMap["redirectTo"]; ok {
		if str, okay := redirectValue.(string); okay {
			//Do a quick website reg ex check?s
			return str
		}
		return ""
	}
	return ""
}
