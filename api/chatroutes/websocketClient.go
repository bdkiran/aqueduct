package chatapi

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

//This code only works for one channel
//This needs to work for more channels

//Apparently this is are not concurrent safe, need to replace with a safer ds
//clients datastructure should only be two way
//[{user, websocket.conn},{user, websocket.conn}]

var clients = make(map[*websocket.Conn]bool) //conected clients
var broadcast = make(chan Message)

var upgrader = websocket.Upgrader{
	//This should not be accepting all orgins, something to think about in the future
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

//Message type
type Message struct {
	//Email    string `json:"email"`
	Username string `json:"username"`
	Message  string `json:"message"`
}

//RegisterConnection registers the websocket connection to our router
func RegisterConnection(r *mux.Router) {
	//Kick off our goRoutine to handle incoming messages
	go handleMessages()
	//handle all incoming connections
	r.HandleFunc("/ws", handleConnections)
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connection made")
	defer ws.Close()

	clients[ws] = true

	for {
		var msg Message
		//Read in the new message as a json and map it message object
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("Error: %v", err)
			delete(clients, ws)
			break
		}
		//send the recieved message to broadcast
		broadcast <- msg
	}
}

func handleMessages() {
	for {
		//grab the message from the broadcast channel
		msg := <-broadcast
		//send it out to every client that is connected
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("Error : %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}
