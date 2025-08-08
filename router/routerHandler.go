package router

import (
	"fmt"
	"log"
	"net/http"
	"server/auth"
	"server/utils"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func handler(w http.ResponseWriter, r *http.Request) {
	tkn := r.URL.Query().Get("token")
	//log.Println(tkn)
	remoteAddr := r.RemoteAddr
	username := WSSGateKeeper(tkn)

	if username == "" || username == "error" {
		log.Println(utils.Yellow("error by: ", remoteAddr, " ", tkn))
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(utils.Red(err))
		return
	}
	utils.LiveConn[conn] = true
	utils.ConnUser[conn] = username
	utils.UserConn[username] = conn
	log.Println("New WebSocket Connection from IP:", remoteAddr, utils.Yellow(" username: ", username), utils.Cyan(" Total Connections : ", CountConn()))
	defer conn.Close()
	HandleConn(conn, username)

}

func RouterHandler() {
	router := mux.NewRouter()
	router.HandleFunc("/wss/chat", handler)
	router.HandleFunc("/login", auth.Login).Methods("POST")
	router.HandleFunc("/register", auth.Register).Methods("POST")
	router.HandleFunc("/validate", auth.GateKeeper).Methods("GET")
	fmt.Println("Server running on", ":8000")
	err := http.ListenAndServe(":8000", router)
	if err != nil {
		fmt.Println("Server Failed to start!")
		panic(err)
	}

}
