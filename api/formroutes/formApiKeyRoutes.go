package formroutes

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	"net/http"

	"github.com/bdkiran/aqueduct/api/response"
)

//Some user identifier as well
type apiKey struct {
	APIKey string `json:"apiKey"`
}

var testKeys = []apiKey{}

//RequestKey Generates and sends back a new form api key.
func RequestKey(w http.ResponseWriter, r *http.Request) {
	generatedKey := generateKey()
	key := apiKey{
		APIKey: generatedKey,
	}
	//Store the key, storing locally for testing purposes
	testKeys = append(testKeys, key)
	response.SendSuccessResponse(w, key)
}

//Random 32 char string, might need a more crytographic secure key
func generateKey() string {
	byteArray := make([]byte, 10)
	if _, err := rand.Read(byteArray); err != nil {
		log.Println("Error generating API key")
	}
	encodedString := hex.EncodeToString(byteArray)

	log.Println(len(encodedString))
	return encodedString
}

//Will need to be a database function. For testing purposes.
func verifyKey(key string) bool {
	for _, validAPIKey := range testKeys {
		if validAPIKey.APIKey == key {
			return true
		}
	}
	return false
}
