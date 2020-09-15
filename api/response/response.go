package response

import (
	"encoding/json"
	"log"
	"net/http"
)

//Status Constants
const (
	StatusSuccess = "success"
	StatusFail    = "fail"
	StatusError   = "error"
)

//ResponseBody structure
type responseBody struct {
	Status  string      `json:"status"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
}

func successResponseBody(data interface{}) responseBody {
	return responseBody{
		Status: StatusSuccess,
		Data:   data,
	}
}

func failResponseBody(data interface{}) responseBody {
	return responseBody{
		Status: StatusFail,
		Data:   data,
	}
}

func errorResponseBody(message string, data interface{}) responseBody {
	return responseBody{
		Status:  StatusError,
		Data:    data,
		Message: message,
	}
}

//SendSuccessResponse a 200 response code
func SendSuccessResponse(w http.ResponseWriter, data interface{}) {
	body := successResponseBody(data)
	b, marshalError := json.Marshal(body)
	//If there is a marshalling error, it is due to the interface passed in,
	//This should only occur if there is a bug with our struct that is passed in.
	if marshalError != nil {
		log.Println(marshalError)
		errorBytes, _ := json.Marshal(errorResponseBody("Error marshaling response", nil))
		writeResponse(w, errorBytes, http.StatusInternalServerError)
	}
	writeResponse(w, b, http.StatusOK)
}

//SendFailResponse a 400 response code
func SendFailResponse(w http.ResponseWriter, data interface{}) {
	body := failResponseBody(data)
	b, marshalError := json.Marshal(body)
	//If there is a marshalling error, it is due to the interface passed in,
	//This should only occur if there is a bug with our struct that is passed in.
	if marshalError != nil {
		log.Println(marshalError)
		errorBytes, _ := json.Marshal(errorResponseBody("Error marshaling response", nil))
		writeResponse(w, errorBytes, http.StatusAccepted)
	}
	writeResponse(w, b, http.StatusBadRequest)
}

//SendErrorResponse sends a 500 response
func SendErrorResponse(w http.ResponseWriter, data interface{}, errorMessage string) {
	body := errorResponseBody(errorMessage, data)
	b, marshalError := json.Marshal(body)
	//If there is a marshalling error, it is due to the interface passed in,
	//This should only occur if there is a bug with our struct that is passed in.
	if marshalError != nil {
		log.Println(marshalError)
		errorBytes, _ := json.Marshal(errorResponseBody("Error marshaling response", nil))
		writeResponse(w, errorBytes, http.StatusInternalServerError)
	}
	writeResponse(w, b, http.StatusInternalServerError)
}

func writeResponse(w http.ResponseWriter, payloadData []byte, statusCode int) {
	w.Header().Set("Content-type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	w.Write(payloadData)
}
