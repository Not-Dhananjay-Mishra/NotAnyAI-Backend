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

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Replace with your frontend port if needed
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		// Handle preflight request
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

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
	corsWrappedRouter := corsMiddleware(router)
	fmt.Println("Server running on", ":8000")
	err := http.ListenAndServe(":8000", corsWrappedRouter)
	if err != nil {
		fmt.Println("Server Failed to start!")
		panic(err)
	}

}
