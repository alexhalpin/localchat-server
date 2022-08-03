package main

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

func outboxRoutine() {
	for {
		select {
		case pkg := <- broadcastChannel:
			broadcastMessage(pkg)
		case msg := <- radarChannel:
			sendRadarPackage(msg)
		}
	}
}

type UserData struct {
	LoginTime string `json:"loginTime"`
	Username string `json:"username"`
	Color string `json:"color"`
	Lat string `json:"lat"`
	Long string `json:"long"`
}

type Client struct {
	UUID string
	Socket *websocket.Conn
	Data UserData
}
var clientMap = make(map[string]*Client)		// map uuid to client struct

// incoming message types
type IncomingMessage struct {
	Kind string `json:"kind"`
	UUID string `json:"uuid"`
	Data map[string]string `json:"data"`
}

//get rid of this
type UserUpdateMessage struct {
	Kind string `json:"kind"`
	UUID string `json:"uuid"`
	Data UserData `json:"data"`
}

var upgrader = websocket.Upgrader{}

func wsHandler(w http.ResponseWriter, r *http.Request){
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	log.Println("WS Endpoint Hit")

	// establish websocket connection from initial request at "/ws"
	conn, err := upgrader.Upgrade(w, r, nil)					
	if err != nil{
		log.Println("Upgrade FAILED: ", err)
		return
	}
	// succesfully connected to client

	// create new  client object and give it a uuid
	var newClient Client
	newClient.UUID = uuid.NewString()
	newClient.Socket = conn

	//login handshake
	userdata, err := loginCheck(newClient)
	if err != nil { 
		log.Println("Login FAILED: ", err)
		return 
	}

	//sucessful login, add userdata from login to map
	newClient.Data.LoginTime = strconv.FormatInt(time.Now().UnixMilli(), 10)
	newClient.Data.Username = (*userdata)["username"]
	newClient.Data.Color = (*userdata)["color"]
	newClient.Data.Lat = (*userdata)["lat"]
	newClient.Data.Long = (*userdata)["long"]

	clientMap[newClient.UUID] = &newClient

	// when this handler returns, remove the client from the client map and close the websocket conn
	defer func() {
		log.Printf("Killing %v", newClient.UUID)
		delete(clientMap, newClient.UUID)
		conn.Close()
	}()

	// Read Message Loop
	for {
		var icMsg IncomingMessage
		conn.SetReadDeadline(time.Now().Add(25 * time.Minute))
		err := conn.ReadJSON(&icMsg)
		if err != nil{
			log.Println("Read FAILED: ", err)
			return
		}

		switch icMsg.Kind {
		case "userUpdate":
			clientMap[icMsg.UUID].Data.Username = icMsg.Data["username"]
			clientMap[icMsg.UUID].Data.Color = icMsg.Data["color"]
			clientMap[icMsg.UUID].Data.Lat = icMsg.Data["lat"]
			clientMap[icMsg.UUID].Data.Long = icMsg.Data["long"]

		case "chat":
			// uuid from incoming -> client map -> username
			// send out message with  username and message

			ogMsg := OutgoingChatMessage {
				Kind: "chat",
				Username: clientMap[icMsg.UUID].Data.Username,
				Color: clientMap[icMsg.UUID].Data.Color,
				Message: icMsg.Data["message"],
			}

			ogPkg := OutgoingPackage {
				SenderUUID: icMsg.UUID,
				Message: ogMsg,
			}

			log.Printf("$ %v: %v", ogMsg.Username, ogMsg.Message)

			broadcastChannel <- ogPkg
		}
	}
}