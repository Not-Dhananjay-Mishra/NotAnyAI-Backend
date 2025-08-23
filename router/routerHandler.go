package router

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"server/auth"
	codingmodel "server/models/CodingModel"
	"server/utils"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Change origin to your frontend domain in production
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

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
		return true // allow all origins for WebSocket
	},
}

func handler(w http.ResponseWriter, r *http.Request) {
	tkn := r.URL.Query().Get("token")
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
func ping(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "up"})
}
func RouterHandler() {
	router := mux.NewRouter()
	router.HandleFunc("/wss/chat", handler)
	router.HandleFunc("/login", auth.Login).Methods("POST")
	router.HandleFunc("/register", auth.Register).Methods("POST")
	router.HandleFunc("/validate", auth.GateKeeper).Methods("GET")
	router.HandleFunc("/pingpong", ping).Methods("GET")
	router.HandleFunc("/api/code", codingmodel.GetRequest).Methods("POST")

	corsWrappedRouter := corsMiddleware(router)

	// Use dynamic port in production
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000" // default for local development
	}
	fmt.Println("Server running on port", port)
	err := http.ListenAndServe(":"+port, corsWrappedRouter)
	if err != nil {
		fmt.Println("Server Failed to start!")
		panic(err)
	}
}
