package router

import (
	"log"
	"server/models"
	"server/utils"

	"github.com/gorilla/websocket"
)

type recived struct {
	Agent string `json:"agent"`
	Query string `json:"query"`
}
type response struct {
	Text string `json:"text"`
}

func HandleConn(conn *websocket.Conn, username string) {
	defer func() {
		log.Println(utils.Magenta("Cleaning up user: ", username))
		delete(utils.LiveConn, conn)
		delete(utils.UserConn, username)
		delete(utils.ConnUser, conn)
		delete(utils.MemoryStore, username)
		conn.Close()
	}()
	client := models.GeminiModel()
	for {
		var data recived
		err := conn.ReadJSON(&data)
		if err != nil {
			log.Println("error reading data")
			return
		}

		receivedData := data
		log.Println(utils.Blue(receivedData.Query))
		if receivedData.Query != "" {

			prompt := receivedData.Query
			models.AddToMemoryUSER(username, prompt)
			aires := models.ModelWithTools(client, utils.MemoryStore[username], username)
			conn.WriteJSON(response{Text: aires})
		}
	}

}
