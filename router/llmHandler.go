package router

import (
	"encoding/base64"
	"log"
	"os"
	"server/models"
	"server/utils"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type recived struct {
	Agent string `json:"agent"`
	Query string `json:"query"`
	Img   string `json:"img"`
}

func HandleConn(conn *websocket.Conn, username string) {
	defer func() {
		log.Println(utils.Magenta("Cleaning up user: ", username))
		//conn.ReadJSON(utils.Response{Text: "Error Encountered Report to Developer"})
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
		if data.Agent == "normal" {
			log.Println(utils.Green("Json"))
		} else if data.Img != "" {
			imgBytes, _ := base64.StdEncoding.DecodeString(data.Img)
			sus := uuid.New().String()
			fileName := "uploads/" + sus + "-" + username + ".png"
			err = os.WriteFile(fileName, imgBytes, 0644)
			if err != nil {
				log.Println("File write error:", err)
			}
			userPrompt := data.Query
			log.Println(utils.Green("Image received and saved as", fileName))
			models.AddImgToMemoryUSER(username, userPrompt, fileName, imgBytes)
			//fmt.Println(utils.MemoryStore[username])
			conn.WriteJSON(utils.Response{Text: "Image received successfully"})
			//models.PrintMemobyUsername(username)
			aires := models.ImageModel(client, fileName, userPrompt, "png", conn, username, imgBytes)
			//aires := models.ModelWithTools(client, utils.MemoryStore[username], username, conn)
			conn.WriteJSON(utils.Response{Text: aires})
			continue
		}

		receivedData := data
		log.Println(utils.Blue(receivedData.Query))
		if receivedData.Query != "" {

			prompt := receivedData.Query
			models.AddToMemoryUSER(username, prompt)
			aires := models.ModelWithTools(client, utils.MemoryStore[username], username, conn)
			conn.WriteJSON(utils.Response{Text: aires})
		}
	}

}
